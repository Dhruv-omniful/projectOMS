package model

import "time"

type Tenant struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID   string    `gorm:"size:100;not null;unique" json:"tenant_id"`
	TenantName string    `gorm:"size:255" json:"tenant_name"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
