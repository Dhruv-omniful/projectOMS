package client

import (
	"context"

	"github.com/omniful/go_commons/log"
)

type IMSClient struct{}

func NewIMSClient() *IMSClient {
	return &IMSClient{}
}

func (c *IMSClient) CheckSKU(ctx context.Context, skuID string) bool {
	// Dummy check: always return true (or put your logic if needed)
	if skuID == "" {
		log.DefaultLogger().Warnf("⚠️ Empty SKU ID provided")
		return false
	}
	log.DefaultLogger().Debugf("✅ Dummy IMS CheckSKU passed for %s", skuID)
	return true
}

func (c *IMSClient) CheckHub(ctx context.Context, hubID string) bool {
	//  Dummy check: always return true (or put your logic if needed)
	if hubID == "" {
		log.DefaultLogger().Warnf("⚠️ Empty Hub ID provided")
		return false
	}
	log.DefaultLogger().Debugf("✅ Dummy IMS CheckHub passed for %s", hubID)
	return true
}
