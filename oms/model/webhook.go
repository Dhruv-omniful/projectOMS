package model

import (
	"time"
)

// Webhook represents a registered webhook for a tenant
type Webhook struct {
	ID          string            `bson:"_id,omitempty" json:"id"`           // MongoDB ID or UUID
	TenantID    string            `bson:"tenant_id" json:"tenant_id"`        // The tenant this webhook belongs to
	CallbackURL string            `bson:"callback_url" json:"callback_url"`  // The target URL for the webhook
	Events      []string          `bson:"events" json:"events"`              // List of event types (e.g. ["order.created"])
	Headers     map[string]string `bson:"headers,omitempty" json:"headers"`  // Optional custom headers (auth tokens etc.)
	Secret      string            `bson:"secret,omitempty" json:"secret"`    // Secret for signing webhook payloads (optional)
	IsActive    bool              `bson:"is_active" json:"is_active"`        // Is the webhook active?
	CreatedAt   time.Time         `bson:"created_at" json:"created_at"`      // Timestamp of creation
	UpdatedAt   time.Time         `bson:"updated_at" json:"updated_at"`      // Timestamp of last update
}
