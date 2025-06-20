package api

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes attaches OMS routes to the engine
func RegisterRoutes(r *gin.Engine, h *Handlers) {
	r.POST("/orders/csv", h.CreateBulkOrder)
	r.POST("/orders/upload-local", h.UploadLocalCSVs)
	// Future routes:
	// r.GET("/orders", h.ListOrders)
	// r.POST("/webhooks", h.RegisterWebhook)
}
