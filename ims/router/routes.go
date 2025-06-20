package routes

import (
	"github.com/gin-gonic/gin"
	"ims/controllers"
)

func RegisterRoutes(r *gin.Engine) {
	// --- Tenants ---
	r.POST("/tenants", controllers.CreateTenant)
	r.GET("/tenants/:id", controllers.GetTenant)
	r.PUT("/tenants/:id", controllers.UpdateTenant)
	r.DELETE("/tenants/:id", controllers.DeleteTenant)
	r.GET("/tenants", controllers.ListTenants)

	// --- Sellers ---
	r.POST("/sellers", controllers.CreateSeller)
	r.GET("/sellers/:id", controllers.GetSeller)
	r.PUT("/sellers/:id", controllers.UpdateSeller)
	r.DELETE("/sellers/:id", controllers.DeleteSeller)
	r.GET("/sellers", controllers.ListSellers)

	// --- Hubs ---
	r.POST("/hubs", controllers.CreateHub)
	r.GET("/hubs/:id", controllers.GetHub)
	r.PUT("/hubs/:id", controllers.UpdateHub)
	r.DELETE("/hubs/:id", controllers.DeleteHub)
	r.GET("/hubs", controllers.ListHubs)

	// --- SKUs ---
	r.POST("/skus", controllers.CreateSKU)
	r.GET("/skus/:id", controllers.GetSKU)
	r.PUT("/skus/:id", controllers.UpdateSKU)
	r.DELETE("/skus/:id", controllers.DeleteSKU)
	r.GET("/skus", controllers.ListSKUs)

	// --- Inventory ---
	r.POST("/inventory", controllers.CreateInventory)
	r.GET("/inventory/:id", controllers.GetInventory)
	r.PUT("/inventory/:id", controllers.UpdateInventory)
	r.DELETE("/inventory/:id", controllers.DeleteInventory)
	r.GET("/inventory", controllers.ListInventory)

	// --- Webhooks ---
	r.POST("/webhooks", controllers.CreateWebhook)
	r.GET("/webhooks/:id", controllers.GetWebhook)
	r.PUT("/webhooks/:id", controllers.UpdateWebhook)
	r.DELETE("/webhooks/:id", controllers.DeleteWebhook)
	r.GET("/webhooks", controllers.ListWebhooks)
}
