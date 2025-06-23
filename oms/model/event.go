package model
import(
    "time"
)
type CreateBulkOrderEvent struct {
	TenantID string `json:"tenant_id"`
	S3Path   string `json:"s3_path"`
	UploadedAt string `json:"uploaded_at"`  // optional timestamp
}


type Order struct {
    ID        string    `bson:"_id,omitempty"`
    TenantID  string    `bson:"tenant_id"`
    SellerID  string    `bson:"seller_id"`
    HubID     string    `bson:"hub_id"`
    SKUID     string    `bson:"sku_id"`
    Quantity  int64     `bson:"quantity"`
    Status    string    `bson:"status"`
    CreatedAt time.Time `bson:"created_at"`
}

type OrderCreated struct {
	OrderID   string    `json:"order_id"`
	TenantID  string    `json:"tenant_id"`
	SellerID  string    `json:"seller_id"`
	HubID     string    `json:"hub_id"`
	SKUID     string    `json:"sku_id"`
	Quantity  int64     `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}