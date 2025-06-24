package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/omniful/go_commons/log"
)

type IMSClient struct {
	BaseURL string
}

func NewIMSClient(baseURL string) *IMSClient {
	return &IMSClient{BaseURL: baseURL}
}

// CheckSKU validates SKU by sku_code
func (c *IMSClient) CheckSKU(ctx context.Context, skuCode string) bool {
	url := fmt.Sprintf("%s/skus/code/%s", c.BaseURL, skuCode)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.DefaultLogger().Errorf(" CheckSKU request error: %v", err)
		return false
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.DefaultLogger().Errorf(" CheckSKU HTTP error: %v", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// CheckHub validates Hub by hub_code
func (c *IMSClient) CheckHub(ctx context.Context, hubCode string) bool {
	url := fmt.Sprintf("%s/hubs/code/%s", c.BaseURL, hubCode)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.DefaultLogger().Errorf(" CheckHub request error: %v", err)
		return false
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.DefaultLogger().Errorf(" CheckHub HTTP error: %v", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
