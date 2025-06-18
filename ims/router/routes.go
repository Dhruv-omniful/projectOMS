package router

import (
	"context"
	"ims/controllers"

	"github.com/omniful/go_commons/http"
	
)

// Initialize sets up all API routes following REST conventions
func Initialize(ctx context.Context, s *http.Server) (err error) {
	apiV1 := s.Engine.Group("/api/v1")

	// Hub routes
	apiV1.POST("/hubs", controllers.CreateHub)
	apiV1.GET("/hubs", controllers.GetAllHubs)
	apiV1.GET("/hubs/:id", controllers.GetHub)
	apiV1.PUT("/hubs/:id", controllers.UpdateHub)
	apiV1.DELETE("/hubs/:id", controllers.DeleteHub)

	// Inventory routes
	apiV1.POST("/inventories", controllers.CreateInventory)
	apiV1.GET("/inventories", controllers.GetAllInventories)
	apiV1.GET("/inventories/:id", controllers.GetInventory)
	apiV1.PUT("/inventories/:id", controllers.UpdateInventory)
	apiV1.DELETE("/inventories/:id", controllers.DeleteInventory)
	apiV1.GET("/inventories/validate/:id", controllers.CheckAndDecrementInventory)

	// Seller routes
	apiV1.POST("/sellers", controllers.CreateSeller)
	apiV1.GET("/sellers", controllers.GetAllSellers)
	apiV1.GET("/sellers/:id", controllers.GetSeller)
	apiV1.PUT("/sellers/:id", controllers.UpdateSeller)
	apiV1.DELETE("/sellers/:id", controllers.DeleteSeller)

	// SKU routes
	apiV1.POST("/skus", controllers.CreateSKU)
	apiV1.GET("/skus", controllers.GetAllSKUs)
	apiV1.GET("/skus/:id", controllers.GetSKU)
	apiV1.PUT("/skus/:id", controllers.UpdateSKU)
	apiV1.DELETE("/skus/:id", controllers.DeleteSKU)
	apiV1.GET("/skus/byTenant/:id", controllers.FetchSKUsByTenant)
	apiV1.GET("/skus/byHub/:id", controllers.FetchSKUsInHub)
	apiV1.GET("/skus/validate/:id", controllers.ValidateSKU)

	// Tenant routes
	apiV1.POST("/tenants", controllers.CreateTenant)
	apiV1.GET("/tenants", controllers.GetAllTenants)
	apiV1.GET("/tenants/:id", controllers.GetTenant)
	apiV1.PUT("/tenants/:id", controllers.UpdateTenant)
	apiV1.DELETE("/tenants/:id", controllers.DeleteTenant)

	// Webhook registration routes
	apiV1.POST("/webhooks/register", controllers.RegisterWebhook)
	apiV1.GET("/webhooks", controllers.GetAllWebhooks)
	apiV1.DELETE("/webhooks/:id", controllers.DeleteWebhook)
	apiV1.POST("/webhooks/trigger/:eventType", controllers.TriggerWebhook) // Optional: simulate webhook trigger

	return nil
}
