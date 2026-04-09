package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/domain"
)

var ErrInvoiceNotFound = errors.New("invoice not found")

type InvoiceRepository struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (r *InvoiceRepository) CreateInvoice(ctx context.Context, invoice domain.Invoice) (domain.Invoice, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Invoice{}, err
	}
	defer tx.Rollback()

	const insertInvoice = `
INSERT INTO notas (status)
VALUES ($1)
RETURNING id, numero, status
`

	var out domain.Invoice
	err = tx.QueryRowContext(ctx, insertInvoice, domain.StatusAberta).Scan(
		&out.ID,
		&out.Numero,
		&out.Status,
	)
	if err != nil {
		return domain.Invoice{}, fmt.Errorf("insert invoice: %w", err)
	}

	const insertItem = `
INSERT INTO nota_itens (nota_id, produto_codigo, quantidade)
VALUES ($1, $2, $3)
RETURNING id, nota_id, produto_codigo, quantidade
`

	out.Itens = make([]domain.InvoiceItem, 0, len(invoice.Itens))
	for _, item := range invoice.Itens {
		var ni domain.InvoiceItem
		err := tx.QueryRowContext(ctx, insertItem, out.ID, item.ProdutoCodigo, item.Quantidade).Scan(
			&ni.ID,
			&ni.NotaID,
			&ni.ProdutoCodigo,
			&ni.Quantidade,
		)
		if err != nil {
			return domain.Invoice{}, fmt.Errorf("insert invoice item: %w", err)
		}
		out.Itens = append(out.Itens, ni)
	}

	if err := tx.Commit(); err != nil {
		return domain.Invoice{}, err
	}

	return out, nil
}

func (r *InvoiceRepository) ListInvoices(ctx context.Context) ([]domain.Invoice, error) {
	const query = `
SELECT id, numero, status, created_at, updated_at
FROM notas
ORDER BY id ASC
`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type invoiceRow struct {
		ID     int64
		Numero int
		Status string
	}

	ids := make([]invoiceRow, 0)
	for rows.Next() {
		var ir invoiceRow
		var createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&ir.ID, &ir.Numero, &ir.Status, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		ids = append(ids, ir)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return []domain.Invoice{}, nil
	}

	invoices := make([]domain.Invoice, 0, len(ids))
	for _, ir := range ids {
		items, err := r.listItemsByInvoice(ctx, ir.ID)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, domain.Invoice{
			ID:     ir.ID,
			Numero: ir.Numero,
			Status: ir.Status,
			Itens:  items,
		})
	}

	return invoices, nil
}

func (r *InvoiceRepository) listItemsByInvoice(ctx context.Context, invoiceID int64) ([]domain.InvoiceItem, error) {
	const query = `
SELECT id, nota_id, produto_codigo, quantidade
FROM nota_itens
WHERE nota_id = $1
ORDER BY id ASC
`

	rows, err := r.db.QueryContext(ctx, query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.InvoiceItem, 0)
	for rows.Next() {
		var item domain.InvoiceItem
		if err := rows.Scan(&item.ID, &item.NotaID, &item.ProdutoCodigo, &item.Quantidade); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *InvoiceRepository) GetInvoiceByID(ctx context.Context, id int64) (domain.Invoice, error) {
	const query = `
SELECT id, numero, status
FROM notas
WHERE id = $1
`

	var inv domain.Invoice
	err := r.db.QueryRowContext(ctx, query, id).Scan(&inv.ID, &inv.Numero, &inv.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Invoice{}, fmt.Errorf("%w: id=%d", ErrInvoiceNotFound, id)
		}
		return domain.Invoice{}, err
	}

	items, err := r.listItemsByInvoice(ctx, inv.ID)
	if err != nil {
		return domain.Invoice{}, err
	}
	inv.Itens = items

	return inv, nil
}

func (r *InvoiceRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	const query = `UPDATE notas SET status = $1 WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return fmt.Errorf("%w: id=%d", ErrInvoiceNotFound, id)
	}

	return nil
}
