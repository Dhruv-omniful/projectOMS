package model

import "time"

type Inventory struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID  string    `gorm:"size:100;not null" json:"tenant_id"`
	SellerID  string    `gorm:"size:100;not null" json:"seller_id"`
	HubCode   string    `gorm:"size:100;not null" json:"hub_code"`
	SKUCode   string    `gorm:"size:100;not null" json:"sku_code"`
	Quantity  int64     `gorm:"default:0" json:"quantity"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Inventory) TableName() string {
	return "inventory"
}
