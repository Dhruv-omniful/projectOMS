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
