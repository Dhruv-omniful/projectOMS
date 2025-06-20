package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"

	"ims/model"
	"ims/postgres"
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
