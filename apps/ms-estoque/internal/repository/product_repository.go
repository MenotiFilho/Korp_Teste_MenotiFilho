package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/domain"
)

var ErrProductCodigoAlreadyExists = errors.New("product codigo already exists")

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, p domain.Product) (domain.Product, error) {
	const query = `
INSERT INTO produtos (codigo, descricao, saldo)
VALUES ($1, $2, $3)
RETURNING id, codigo, descricao, saldo
`

	var out domain.Product
	err := r.db.QueryRowContext(ctx, query, p.Codigo, p.Descricao, p.Saldo).Scan(
		&out.ID,
		&out.Codigo,
		&out.Descricao,
		&out.Saldo,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.Product{}, fmt.Errorf("%w: codigo=%s", ErrProductCodigoAlreadyExists, p.Codigo)
		}
		return domain.Product{}, err
	}

	return out, nil
}

func (r *ProductRepository) ListProducts(ctx context.Context) ([]domain.Product, error) {
	const query = `
SELECT id, codigo, descricao, saldo
FROM produtos
ORDER BY id ASC
`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]domain.Product, 0)
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.Codigo, &p.Descricao, &p.Saldo); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
