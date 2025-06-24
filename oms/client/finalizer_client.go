package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/omniful/go_commons/log"
	"go.mongodb.org/mongo-driver/bson"
)

// IMSInventory defines the expected response structure from IMS
type IMSInventory struct {
	ID       int64  `json:"id"`
	TenantID string `json:"tenant_id"`
	SellerID string `json:"seller_id"`
	HubCode  string `json:"hub_code"`
	SKUCode  string `json:"sku_code"`
	Quantity int64  `json:"quantity"`
}

// FetchInventory calls IMS to get inventory info
func FetchInventory(ctx context.Context, baseURL, tenantID, sellerID, hubCode, skuCode string) (*IMSInventory, error) {
	logger := log.DefaultLogger()

	// Build query URL safely using net/url
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid baseURL: %w", err)
	}
	u.Path = "/inventory/query"

	// Add query parameters
	q := u.Query()
	q.Set("tenant_id", tenantID)
	q.Set("seller_id", sellerID)
	q.Set("hub_code", hubCode)
	q.Set("sku_code", skuCode)
	u.RawQuery = q.Encode()

	finalURL := u.String()
	logger.Infof("IMS request URL: %s", finalURL)

	// Build and make request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, finalURL, nil)
	if err != nil {
		return nil, fmt.Errorf("new request error: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %w", err)
	}
	defer resp.Body.Close()

	logger.Infof("IMS response status: %s", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IMS returned status %d", resp.StatusCode)
	}

	var inv IMSInventory
	if err := json.NewDecoder(resp.Body).Decode(&inv); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	return &inv, nil
}

// ConsumeInventory calls IMS to consume stock
func ConsumeInventory(ctx context.Context, baseURL, tenantID, sellerID, hubCode, skuCode string, qty int64) error {
	url := fmt.Sprintf("%s/inventory/consume", baseURL)
	payload := map[string]interface{}{
		"tenant_id": tenantID,
		"seller_id": sellerID,
		"hub_code":  hubCode,
		"sku_code":  skuCode,
		"quantity":  qty,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("new request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("IMS returned status %d", resp.StatusCode)
	}

	return nil
}

// UpdateOrderStatus calls OMS to update the order's status
type UpdateOrderStatusRequest struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

func UpdateOrderStatus(ctx context.Context, _ string, req UpdateOrderStatusRequest) error {
	logger := log.DefaultLogger()

	coll, err := GetOrdersCollection(ctx)
	if err != nil {
		return fmt.Errorf("get collection error: %w", err)
	}

	filter := bson.M{"_id": req.OrderID}
	update := bson.M{"$set": bson.M{"status": req.Status}}

	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Errorf(" Failed to update order status: %v", err)
		return fmt.Errorf("update error: %w", err)
	}

	if result.MatchedCount == 0 {
		logger.Warnf(" No order found with ID %s to update", req.OrderID)
		return fmt.Errorf("no order found with ID %s", req.OrderID)
	}

	logger.Infof(" Updated order status: OrderID=%s Status=%s", req.OrderID, req.Status)
	return nil
}
