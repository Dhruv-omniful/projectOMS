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

// CreateTenant handles POST /tenants
func CreateTenant(c *gin.Context) {
	var tenant model.Tenant
	if err := c.ShouldBindJSON(&tenant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}

	now := time.Now().UTC()
	tenant.CreatedAt = now
	tenant.UpdatedAt = now

	db := pr.DB.GetMasterDB(c.Request.Context())
	if err := db.Create(&tenant).Error; err != nil {
		log.DefaultLogger().Errorf("CreateTenant DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.create_tenant_failed")})
		return
	}

	c.JSON(http.StatusCreated, tenant)
}

// GetTenant handles GET /tenants/:id
func GetTenant(c *gin.Context) {
	id := c.Param("id")
	var tenant model.Tenant

	db := pr.DB.GetSlaveDB(c.Request.Context())
	if err := db.First(&tenant, "id = ?", id).Error; err != nil {
		log.DefaultLogger().Errorf("GetTenant DB error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.Translate(c, "error.tenant_not_found")})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

// UpdateTenant handles PUT /tenants/:id
func UpdateTenant(c *gin.Context) {
	id := c.Param("id")
	var tenant model.Tenant

	db := pr.DB.GetMasterDB(c.Request.Context())
	if err := db.First(&tenant, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.Translate(c, "error.tenant_not_found")})
		return
	}

	if err := c.ShouldBindJSON(&tenant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}

	tenant.UpdatedAt = time.Now().UTC()

	if err := db.Save(&tenant).Error; err != nil {
		log.DefaultLogger().Errorf("UpdateTenant DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.update_tenant_failed")})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

// DeleteTenant handles DELETE /tenants/:id
func DeleteTenant(c *gin.Context) {
	id := c.Param("id")

	db := pr.DB.GetMasterDB(c.Request.Context())
	if err := db.Delete(&model.Tenant{}, "id = ?", id).Error; err != nil {
		log.DefaultLogger().Errorf("DeleteTenant DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.delete_tenant_failed")})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListTenants handles GET /tenants
func ListTenants(c *gin.Context) {
	var tenants []model.Tenant

	db := pr.DB.GetSlaveDB(c.Request.Context())
	if err := db.Find(&tenants).Error; err != nil {
		log.DefaultLogger().Errorf("ListTenants DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.list_tenants_failed")})
		return
	}

	c.JSON(http.StatusOK, tenants)
}
