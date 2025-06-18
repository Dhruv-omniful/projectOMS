package model

import "time"

type Hub struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID  string    `gorm:"size:100;not null" json:"tenant_id"`
	SellerID  string    `gorm:"size:100;not null" json:"seller_id"`
	HubCode   string    `gorm:"size:100;not null;unique" json:"hub_code"`
	HubName   string    `gorm:"size:255" json:"hub_name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
