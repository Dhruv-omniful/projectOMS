package api

import (
	"net/http"

	"github.com/dhruv/oms/client"
	"github.com/dhruv/oms/model"
	service "github.com/dhruv/oms/services"
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/log"
)

// Handlers wraps services
type Handlers struct {
	OrderService *service.OrderService
}

// NewHandlers creates handlers with dependencies
func NewHandlers(orderService *service.OrderService) *Handlers {
	return &Handlers{
		OrderService: orderService,
	}
}

// CreateBulkOrder handles POST /orders/csv
func (h *Handlers) CreateBulkOrder(c *gin.Context) {
	var req struct {
		Path string `json:"path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf(" Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: missing or bad path field",
		})
		return
	}

	ctx := c.Request.Context()
	if err := h.OrderService.ProcessCSV(ctx, req.Path); err != nil {
		log.Errorf(" Failed to process CSV: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process CSV file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "CSV file processed successfully (S3 path validated)",
	})
}

// RegisterWebhook handles POST /webhooks
func (h *Handlers) RegisterWebhook(c *gin.Context) {
	var req model.Webhook
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf(" Invalid webhook request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}
	if err := client.SaveWebhook(c.Request.Context(), &req); err != nil {
		log.Errorf(" Failed to save webhook: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save webhook"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Webhook registered"})
}
