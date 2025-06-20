package model

type CreateBulkOrderEvent struct {
	TenantID string `json:"tenant_id"`
	S3Path   string `json:"s3_path"`
	UploadedAt string `json:"uploaded_at"`  // optional timestamp
}
