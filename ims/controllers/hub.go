package controllers

import (
	"net/http"
	"time"

	"encoding/json"
	"ims/model"
	pr "ims/postgres"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
)

// CreateHub handles POST /hubs
func CreateHub(c *gin.Context) {
	var hub model.Hub
	if err := c.ShouldBindJSON(&hub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}

	now := time.Now().UTC()
	hub.CreatedAt = now
	hub.UpdatedAt = now

	db := pr.DB.GetMasterDB(c.Request.Context())
	if err := db.Create(&hub).Error; err != nil {
		log.DefaultLogger().Errorf("CreateHub DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.create_hub_failed")})
		return
	}

	c.JSON(http.StatusCreated, hub)
}

// GetHub handles GET /hubs/:id
func GetHub(c *gin.Context) {
	id := c.Param("id")

	// Try cache first
	if cached, err := pr.RedisClient.Get(c.Request.Context(), "hub:"+id); err == nil {
		c.JSON(http.StatusOK, cached)
		return
	}

	var hub model.Hub

	db := pr.DB.GetSlaveDB(c.Request.Context())
	if err := db.First(&hub, "id = ?", id).Error; err != nil {
		log.DefaultLogger().Errorf("GetHub DB error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.Translate(c, "error.hub_not_found")})
		return
	}

	if b, err := json.Marshal(hub); err == nil {
		// ignore both return values
		_, _ = pr.RedisClient.Set(c.Request.Context(), "hub:"+id, string(b), 5*time.Minute)
	}

	c.JSON(http.StatusOK, hub)
}

// UpdateHub handles PUT /hubs/:id
func UpdateHub(c *gin.Context) {
	id := c.Param("id")
	var hub model.Hub

	db := pr.DB.GetMasterDB(c.Request.Context())
	if err := db.First(&hub, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.Translate(c, "error.hub_not_found")})
		return
	}

	if err := c.ShouldBindJSON(&hub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}

	hub.UpdatedAt = time.Now().UTC()

	if err := db.Save(&hub).Error; err != nil {
		log.DefaultLogger().Errorf("UpdateHub DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.update_hub_failed")})
		return
	}

	c.JSON(http.StatusOK, hub)
}

// DeleteHub handles DELETE /hubs/:id
func DeleteHub(c *gin.Context) {
	id := c.Param("id")

	db := pr.DB.GetMasterDB(c.Request.Context())
	if err := db.Delete(&model.Hub{}, "id = ?", id).Error; err != nil {
		log.DefaultLogger().Errorf("DeleteHub DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.delete_hub_failed")})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListHubs handles GET /hubs
func ListHubs(c *gin.Context) {
	var hubs []model.Hub

	db := pr.DB.GetSlaveDB(c.Request.Context())
	if err := db.Find(&hubs).Error; err != nil {
		log.DefaultLogger().Errorf("ListHubs DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.list_hubs_failed")})
		return
	}

	c.JSON(http.StatusOK, hubs)
}

// GetHubByCode handles GET /hubs/code/:hub_code
func GetHubByCode(c *gin.Context) {
	hubCode := c.Param("hub_code")
	var hub model.Hub

	db := pr.DB.GetSlaveDB(c.Request.Context())
	if err := db.Where("hub_code = ?", hubCode).First(&hub).Error; err != nil {
		log.DefaultLogger().Errorf("GetHubByCode DB error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.Translate(c, "error.hub_not_found")})
		return
	}

	c.JSON(http.StatusOK, hub)
}
