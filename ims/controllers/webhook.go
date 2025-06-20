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

// CreateWebhook handles POST /webhooks
func CreateWebhook(c *gin.Context) {
	var webhook model.WebhookRegistration
	if err := c.ShouldBindJSON(&webhook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}

	webhook.CreatedAt = time.Now().UTC()

	db := pr.DB.GetMasterDB(c.Request.Context())
	if err := db.Create(&webhook).Error; err != nil {
		log.DefaultLogger().Errorf("CreateWebhook DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.create_webhook_failed")})
		return
	}

	c.JSON(http.StatusCreated, webhook)
}

// GetWebhook handles GET /webhooks/:id
func GetWebhook(c *gin.Context) {
	id := c.Param("id")
	var webhook model.WebhookRegistration

	db := pr.DB.GetSlaveDB(c.Request.Context())
	if err := db.First(&webhook, "id = ?", id).Error; err != nil {
		log.DefaultLogger().Errorf("GetWebhook DB error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.Translate(c, "error.webhook_not_found")})
		return
	}

	c.JSON(http.StatusOK, webhook)
}

// UpdateWebhook handles PUT /webhooks/:id
func UpdateWebhook(c *gin.Context) {
	id := c.Param("id")
	var webhook model.WebhookRegistration

	db := pr.DB.GetMasterDB(c.Request.Context())
	if err := db.First(&webhook, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.Translate(c, "error.webhook_not_found")})
		return
	}

	if err := c.ShouldBindJSON(&webhook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}

	// We don't update CreatedAt
	if err := db.Save(&webhook).Error; err != nil {
		log.DefaultLogger().Errorf("UpdateWebhook DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.update_webhook_failed")})
		return
	}

	c.JSON(http.StatusOK, webhook)
}

// DeleteWebhook handles DELETE /webhooks/:id
func DeleteWebhook(c *gin.Context) {
	id := c.Param("id")

	db := pr.DB.GetMasterDB(c.Request.Context())
	if err := db.Delete(&model.WebhookRegistration{}, "id = ?", id).Error; err != nil {
		log.DefaultLogger().Errorf("DeleteWebhook DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.delete_webhook_failed")})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListWebhooks handles GET /webhooks
func ListWebhooks(c *gin.Context) {
	var webhooks []model.WebhookRegistration

	db := pr.DB.GetSlaveDB(c.Request.Context())
	if err := db.Find(&webhooks).Error; err != nil {
		log.DefaultLogger().Errorf("ListWebhooks DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.list_webhooks_failed")})
		return
	}

	c.JSON(http.StatusOK, webhooks)
}
