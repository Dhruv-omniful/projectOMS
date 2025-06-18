package model

import "time"

type Seller struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID  string    `gorm:"size:100;not null" json:"tenant_id"`
	SellerID  string    `gorm:"size:100;not null;unique" json:"seller_id"`
	SellerName string   `gorm:"size:255" json:"seller_name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
