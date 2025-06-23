package controllers

import (
	"net/http"
	"time"

	"ims/model"
	pr "ims/postgres"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
	"gorm.io/gorm/clause"
)

// CreateInventory handles POST /inventory (upsert logic can be added if needed)
func CreateInventory(c *gin.Context) {
	var inventory model.Inventory
	if err := c.ShouldBindJSON(&inventory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}

	inventory.UpdatedAt = time.Now().UTC()

	db := pr.DB.GetMasterDB(c.Request.Context())
	if err := db.Create(&inventory).Error; err != nil {
		log.DefaultLogger().Errorf("CreateInventory DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.create_inventory_failed")})
		return
	}

	c.JSON(http.StatusCreated, inventory)
}

// GetInventory handles GET /inventory/:id
func GetInventory(c *gin.Context) {
	id := c.Param("id")
	var inventory model.Inventory

	db := pr.DB.GetSlaveDB(c.Request.Context())
	if err := db.First(&inventory, "id = ?", id).Error; err != nil {
		log.DefaultLogger().Errorf("GetInventory DB error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.Translate(c, "error.inventory_not_found")})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

// UpdateInventory handles PUT /inventory/:id
func UpdateInventory(c *gin.Context) {
	id := c.Param("id")
	var inventory model.Inventory

	db := pr.DB.GetMasterDB(c.Request.Context())
	if err := db.First(&inventory, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.Translate(c, "error.inventory_not_found")})
		return
	}

	if err := c.ShouldBindJSON(&inventory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}

	inventory.UpdatedAt = time.Now().UTC()

	if err := db.Save(&inventory).Error; err != nil {
		log.DefaultLogger().Errorf("UpdateInventory DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.update_inventory_failed")})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

// DeleteInventory handles DELETE /inventory/:id
func DeleteInventory(c *gin.Context) {
	id := c.Param("id")

	db := pr.DB.GetMasterDB(c.Request.Context())
	if err := db.Delete(&model.Inventory{}, "id = ?", id).Error; err != nil {
		log.DefaultLogger().Errorf("DeleteInventory DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.delete_inventory_failed")})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListInventory handles GET /inventory
func ListInventory(c *gin.Context) {
	var inventories []model.Inventory

	db := pr.DB.GetSlaveDB(c.Request.Context())
	if err := db.Find(&inventories).Error; err != nil {
		log.DefaultLogger().Errorf("ListInventory DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.list_inventory_failed")})
		return
	}

	c.JSON(http.StatusOK, inventories)
}

// QueryInventory handles GET /inventory/query?tenant_id=...&seller_id=...&hub_code=...&sku_code=...
func QueryInventory(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	sellerID := c.Query("seller_id")
	hubCode := c.Query("hub_code")
	skuCode := c.Query("sku_code")

	if tenantID == "" || sellerID == "" || hubCode == "" || skuCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required query params"})
		return
	}

	var inventory model.Inventory
	db := pr.DB.GetSlaveDB(c.Request.Context())
	if err := db.Where("tenant_id = ? AND seller_id = ? AND hub_code = ? AND sku_code = ?", tenantID, sellerID, hubCode, skuCode).
		First(&inventory).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventory not found"})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

// ConsumeInventory handles POST /inventory/consume
func ConsumeInventory(c *gin.Context) {
	var req struct {
		TenantID string `json:"tenant_id"`
		SellerID string `json:"seller_id"`
		HubCode  string `json:"hub_code"`
		SKUCode  string `json:"sku_code"`
		Quantity int64  `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	db := pr.DB.GetMasterDB(c.Request.Context())

	var inventory model.Inventory
	// Lock row for update
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("tenant_id = ? AND seller_id = ? AND hub_code = ? AND sku_code = ?", req.TenantID, req.SellerID, req.HubCode, req.SKUCode).
		First(&inventory).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventory not found"})
		return
	}

	log.Infof("🔍 Fetched inventory before update: %+v", inventory)

	if inventory.Quantity < req.Quantity {
		c.JSON(http.StatusConflict, gin.H{"error": "Insufficient inventory"})
		return
	}

	newQty := inventory.Quantity - req.Quantity
	updatedAt := time.Now().UTC()

	// Use Updates instead of Save for reliability
	if err := db.Model(&inventory).
		Where("id = ?", inventory.ID).
		Updates(map[string]interface{}{
			"quantity":   newQty,
			"updated_at": updatedAt,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory"})
		return
	}

	log.Infof("✅ Inventory updated: ID=%d New Quantity=%d", inventory.ID, newQty)

	c.JSON(http.StatusOK, gin.H{
		"message":   "Inventory consumed",
		"remaining": newQty,
	})
}
