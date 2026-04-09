CREATE TABLE IF NOT EXISTS nota_itens (
    id BIGSERIAL PRIMARY KEY,
    nota_id BIGINT NOT NULL REFERENCES notas(id) ON DELETE CASCADE,
    produto_codigo VARCHAR(64) NOT NULL,
    quantidade INTEGER NOT NULL CHECK (quantidade > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_nota_itens_nota_id ON nota_itens(nota_id);
