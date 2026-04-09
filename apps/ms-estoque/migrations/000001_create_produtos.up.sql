CREATE TABLE IF NOT EXISTS produtos (
    id BIGSERIAL PRIMARY KEY,
    codigo VARCHAR(64) NOT NULL UNIQUE,
    descricao VARCHAR(255) NOT NULL,
    saldo INTEGER NOT NULL CHECK (saldo >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_produtos_set_updated_at
BEFORE UPDATE ON produtos
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
