package model

import "time"

type WebhookRegistration struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	TenantID  string    `gorm:"not null" json:"tenant_id"`
	URL       string    `gorm:"not null" json:"url"`
	EventType string    `gorm:"not null" json:"event_type"` // e.g. "order.created"
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
