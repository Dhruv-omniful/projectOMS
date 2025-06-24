package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/dhruv/oms/model"
	"github.com/omniful/go_commons/log"
)

// NotifyWebhooks sends event data to the webhook URL
func NotifyWebhooks(ctx context.Context, tenantID, eventType string, payload interface{}) {
	webhooks, err := GetWebhooksForTenantAndEvent(ctx, tenantID, eventType)
	if err != nil {
		log.DefaultLogger().Errorf("❌ Failed to fetch webhooks: %v", err)
		return
	}

	if len(webhooks) == 0 {
		log.DefaultLogger().Infof("⚠️ No webhooks registered for tenant=%s event=%s", tenantID, eventType)
		return
	}

	body, err := json.Marshal(map[string]interface{}{
		"tenant_id": tenantID,
		"event":     eventType,
		"data":      payload,
	})
	if err != nil {
		log.DefaultLogger().Errorf("❌ Failed to marshal webhook payload: %v", err)
		return
	}

	for _, wh := range webhooks {
		if !wh.IsActive {
			continue
		}

		go sendWebhookRequest(ctx, wh, body)
	}
}

func sendWebhookRequest(ctx context.Context, wh model.Webhook, body []byte) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, wh.CallbackURL, bytes.NewReader(body))
	if err != nil {
		log.DefaultLogger().Errorf("❌ Failed to create webhook request for %s: %v", wh.CallbackURL, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Add custom headers if present
	for k, v := range wh.Headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.DefaultLogger().Errorf("❌ Webhook POST failed to %s: %v", wh.CallbackURL, err)
		return
	}
	defer resp.Body.Close()

	log.DefaultLogger().Infof("✅ Webhook sent to %s: status=%d", wh.CallbackURL, resp.StatusCode)
}
