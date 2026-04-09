package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/domain"
)

var ErrProductCodigoAlreadyExists = errors.New("product codigo already exists")

var (
	ErrProductNotFound          = errors.New("product not found")
	ErrProductInsufficientStock = errors.New("product insufficient stock")
	ErrInvalidDecreaseItem      = errors.New("invalid stock decrease item")
	ErrIdempotencyKeyRequired   = errors.New("idempotency key is required")
)

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

func (r *ProductRepository) DecreaseStock(ctx context.Context, items []domain.StockDecreaseItem, idempotencyKey string) error {
	if len(items) == 0 {
		return nil
	}

	idempotencyKey = strings.TrimSpace(idempotencyKey)
	if idempotencyKey == "" {
		return ErrIdempotencyKeyRequired
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `INSERT INTO estoque_baixas (idempotency_key) VALUES ($1)`, idempotencyKey); err != nil {
		if isUniqueViolation(err) {
			return nil
		}
		return err
	}

	normalized := make([]domain.StockDecreaseItem, 0, len(items))
	for _, item := range items {
		codigo := strings.TrimSpace(item.Codigo)
		if codigo == "" || item.Quantidade <= 0 {
			return ErrInvalidDecreaseItem
		}
		normalized = append(normalized, domain.StockDecreaseItem{Codigo: codigo, Quantidade: item.Quantidade})
	}

	sort.Slice(normalized, func(i, j int) bool {
		if normalized[i].Codigo == normalized[j].Codigo {
			return normalized[i].Quantidade < normalized[j].Quantidade
		}
		return normalized[i].Codigo < normalized[j].Codigo
	})

	for _, item := range normalized {
		codigo := item.Codigo

		var saldo int
		err := tx.QueryRowContext(ctx, "SELECT saldo FROM produtos WHERE codigo = $1 FOR UPDATE", codigo).Scan(&saldo)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("%w: codigo=%s", ErrProductNotFound, codigo)
			}
			return err
		}

		if saldo < item.Quantidade {
			return fmt.Errorf("%w: codigo=%s saldo=%d solicitado=%d", ErrProductInsufficientStock, codigo, saldo, item.Quantidade)
		}

		if _, err := tx.ExecContext(ctx, "UPDATE produtos SET saldo = saldo - $1 WHERE codigo = $2", item.Quantidade, codigo); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
