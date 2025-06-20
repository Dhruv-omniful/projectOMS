CREATE TABLE IF NOT EXISTS skus (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(100) NOT NULL,
    seller_id VARCHAR(100) NOT NULL,
    sku_code VARCHAR(100) NOT NULL UNIQUE,
    sku_name VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_skus_tenant_id ON skus (tenant_id);
CREATE INDEX IF NOT EXISTS idx_skus_seller_id ON skus (seller_id);