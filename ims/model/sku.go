package model

import "time"

type SKU struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID  string    `gorm:"size:100;not null" json:"tenant_id"`
	SellerID  string    `gorm:"size:100;not null" json:"seller_id"`
	SKUCode   string    `gorm:"size:100;not null;unique" json:"sku_code"`
	SKUName   string    `gorm:"size:255" json:"sku_name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
