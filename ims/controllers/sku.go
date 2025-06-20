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

// CreateSKU handles POST /skus
func CreateSKU(c *gin.Context) {
	var sku model.SKU
	if err := c.ShouldBindJSON(&sku); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}

	now := time.Now().UTC()
	sku.CreatedAt = now
	sku.UpdatedAt = now

	db := pr.DB.GetMasterDB(c.Request.Context())
	if err := db.Create(&sku).Error; err != nil {
		log.DefaultLogger().Errorf("CreateSKU DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.create_sku_failed")})
		return
	}

	c.JSON(http.StatusCreated, sku)
}

// GetSKU handles GET /skus/:id
func GetSKU(c *gin.Context) {
	id := c.Param("id")
	var sku model.SKU

	db := pr.DB.GetSlaveDB(c.Request.Context())
	if err := db.First(&sku, "id = ?", id).Error; err != nil {
		log.DefaultLogger().Errorf("GetSKU DB error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.Translate(c, "error.sku_not_found")})
		return
	}

	c.JSON(http.StatusOK, sku)
}

// UpdateSKU handles PUT /skus/:id
func UpdateSKU(c *gin.Context) {
	id := c.Param("id")
	var sku model.SKU

	db := pr.DB.GetMasterDB(c.Request.Context())
	if err := db.First(&sku, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.Translate(c, "error.sku_not_found")})
		return
	}

	if err := c.ShouldBindJSON(&sku); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}

	sku.UpdatedAt = time.Now().UTC()

	if err := db.Save(&sku).Error; err != nil {
		log.DefaultLogger().Errorf("UpdateSKU DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.update_sku_failed")})
		return
	}

	c.JSON(http.StatusOK, sku)
}

// DeleteSKU handles DELETE /skus/:id
func DeleteSKU(c *gin.Context) {
	id := c.Param("id")

	db := pr.DB.GetMasterDB(c.Request.Context())
	if err := db.Delete(&model.SKU{}, "id = ?", id).Error; err != nil {
		log.DefaultLogger().Errorf("DeleteSKU DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.delete_sku_failed")})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListSKUs handles GET /skus
func ListSKUs(c *gin.Context) {
	var skus []model.SKU

	db := pr.DB.GetSlaveDB(c.Request.Context())
	if err := db.Find(&skus).Error; err != nil {
		log.DefaultLogger().Errorf("ListSKUs DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.list_skus_failed")})
		return
	}

	c.JSON(http.StatusOK, skus)
}
