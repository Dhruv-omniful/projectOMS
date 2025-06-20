CREATE TABLE IF NOT EXISTS inventory (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(100) NOT NULL,
    seller_id VARCHAR(100) NOT NULL,
    hub_code VARCHAR(100) NOT NULL,
    sku_code VARCHAR(100) NOT NULL,
    quantity BIGINT DEFAULT 0,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_inventory_tenant_id ON inventory (tenant_id);
CREATE INDEX IF NOT EXISTS idx_inventory_seller_id ON inventory (seller_id);
CREATE INDEX IF NOT EXISTS idx_inventory_hub_code ON inventory (hub_code);
CREATE INDEX IF NOT EXISTS idx_inventory_sku_code ON inventory (sku_code);