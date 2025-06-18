CREATE TABLE IF NOT EXISTS hubs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(100) NOT NULL,
    seller_id VARCHAR(100) NOT NULL,
    hub_code VARCHAR(100) NOT NULL UNIQUE,
    hub_name VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_hubs_tenant_id ON hubs (tenant_id);
CREATE INDEX IF NOT EXISTS idx_hubs_seller_id ON hubs (seller_id);