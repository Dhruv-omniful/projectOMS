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

// CreateSeller handles POST /sellers
func CreateSeller(c *gin.Context) {
	var seller model.Seller
	if err := c.ShouldBindJSON(&seller); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}

	now := time.Now().UTC()
	seller.CreatedAt = now
	seller.UpdatedAt = now

	db := pr.DB.GetMasterDB(c.Request.Context())
	if err := db.Create(&seller).Error; err != nil {
		log.DefaultLogger().Errorf("CreateSeller DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.create_seller_failed")})
		return
	}

	c.JSON(http.StatusCreated, seller)
}

// GetSeller handles GET /sellers/:id
func GetSeller(c *gin.Context) {
	id := c.Param("id")
	var seller model.Seller

	db := pr.DB.GetSlaveDB(c.Request.Context())
	if err := db.First(&seller, "id = ?", id).Error; err != nil {
		log.DefaultLogger().Errorf("GetSeller DB error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.Translate(c, "error.seller_not_found")})
		return
	}

	c.JSON(http.StatusOK, seller)
}

// UpdateSeller handles PUT /sellers/:id
func UpdateSeller(c *gin.Context) {
	id := c.Param("id")
	var seller model.Seller

	db := pr.DB.GetMasterDB(c.Request.Context())
	if err := db.First(&seller, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.Translate(c, "error.seller_not_found")})
		return
	}

	if err := c.ShouldBindJSON(&seller); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}

	seller.UpdatedAt = time.Now().UTC()

	if err := db.Save(&seller).Error; err != nil {
		log.DefaultLogger().Errorf("UpdateSeller DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.update_seller_failed")})
		return
	}

	c.JSON(http.StatusOK, seller)
}

// DeleteSeller handles DELETE /sellers/:id
func DeleteSeller(c *gin.Context) {
	id := c.Param("id")

	db := pr.DB.GetMasterDB(c.Request.Context())
	if err := db.Delete(&model.Seller{}, "id = ?", id).Error; err != nil {
		log.DefaultLogger().Errorf("DeleteSeller DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.delete_seller_failed")})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListSellers handles GET /sellers
func ListSellers(c *gin.Context) {
	var sellers []model.Seller

	db := pr.DB.GetSlaveDB(c.Request.Context())
	if err := db.Find(&sellers).Error; err != nil {
		log.DefaultLogger().Errorf("ListSellers DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.list_sellers_failed")})
		return
	}

	c.JSON(http.StatusOK, sellers)
}
