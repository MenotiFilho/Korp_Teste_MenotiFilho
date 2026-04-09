package repository

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/domain"
)

const defaultTestDatabaseURL = "postgres://postgres:postgres@localhost:5434/faturamento?sslmode=disable"

func TestInvoiceRepository_CreateAndListInvoices_ShouldPersistAndReturn(t *testing.T) {
	// Arrange
	db := openTestDB(t)
	repo := NewInvoiceRepository(db)
	ctx := context.Background()

	items := []domain.InvoiceItem{
		{ProdutoCodigo: "P-001", Quantidade: 2},
		{ProdutoCodigo: "P-002", Quantidade: 5},
	}
	invoice, err := domain.NewInvoice(items)
	if err != nil {
		t.Fatalf("unexpected domain error: %v", err)
	}

	// Act
	created, err := repo.CreateInvoice(ctx, invoice)
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}

	invoices, err := repo.ListInvoices(ctx)
	if err != nil {
		t.Fatalf("unexpected list error: %v", err)
	}

	// Assert
	if created.ID <= 0 {
		t.Fatalf("expected generated id > 0, got %d", created.ID)
	}
	if created.Numero <= 0 {
		t.Fatalf("expected generated numero > 0, got %d", created.Numero)
	}
	if created.Status != domain.StatusAberta {
		t.Fatalf("expected status %q, got %q", domain.StatusAberta, created.Status)
	}
	if len(created.Itens) != 2 {
		t.Fatalf("expected 2 items, got %d", len(created.Itens))
	}
	if created.Itens[0].ProdutoCodigo != "P-001" || created.Itens[1].ProdutoCodigo != "P-002" {
		t.Fatalf("unexpected items: %+v", created.Itens)
	}

	if len(invoices) != 1 {
		t.Fatalf("expected 1 invoice in list, got %d", len(invoices))
	}
	if invoices[0].Numero != created.Numero {
		t.Fatalf("expected same numero %d, got %d", created.Numero, invoices[0].Numero)
	}
}

func TestInvoiceRepository_SequentialNumbering_ShouldAssignIncrementalNumbers(t *testing.T) {
	// Arrange
	db := openTestDB(t)
	repo := NewInvoiceRepository(db)
	ctx := context.Background()

	items1 := []domain.InvoiceItem{{ProdutoCodigo: "P-001", Quantidade: 1}}
	items2 := []domain.InvoiceItem{{ProdutoCodigo: "P-002", Quantidade: 1}}
	invoice1, _ := domain.NewInvoice(items1)
	invoice2, _ := domain.NewInvoice(items2)

	// Act
	created1, err := repo.CreateInvoice(ctx, invoice1)
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}
	created2, err := repo.CreateInvoice(ctx, invoice2)
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}

	// Assert
	if created2.Numero <= created1.Numero {
		t.Fatalf("expected sequential numbering: %d > %d", created2.Numero, created1.Numero)
	}
}

func TestInvoiceRepository_UpdateStatus_ShouldChangeStatusAndTimestamp(t *testing.T) {
	// Arrange
	db := openTestDB(t)
	repo := NewInvoiceRepository(db)
	ctx := context.Background()

	items := []domain.InvoiceItem{{ProdutoCodigo: "P-001", Quantidade: 1}}
	invoice, _ := domain.NewInvoice(items)
	created, _ := repo.CreateInvoice(ctx, invoice)

	// Act
	err := repo.UpdateStatus(ctx, created.ID, domain.StatusFechada)
	if err != nil {
		t.Fatalf("unexpected update error: %v", err)
	}

	invoices, _ := repo.ListInvoices(ctx)

	// Assert
	if len(invoices) != 1 {
		t.Fatalf("expected 1 invoice, got %d", len(invoices))
	}
	if invoices[0].Status != domain.StatusFechada {
		t.Fatalf("expected status %q, got %q", domain.StatusFechada, invoices[0].Status)
	}
}

func TestInvoiceRepository_UpdateStatus_WhenInvoiceNotFound_ShouldReturnError(t *testing.T) {
	// Arrange
	db := openTestDB(t)
	repo := NewInvoiceRepository(db)
	ctx := context.Background()

	// Act
	err := repo.UpdateStatus(ctx, 99999, domain.StatusFechada)

	// Assert
	if err == nil {
		t.Fatal("expected error for non-existent invoice")
	}
}

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	if testing.Short() {
		t.Skip("integration test skipped in short mode")
	}

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = defaultTestDatabaseURL
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping database (%s): %v", dbURL, err)
	}

	if _, err := db.Exec("TRUNCATE TABLE nota_itens RESTART IDENTITY CASCADE"); err != nil {
		t.Fatalf("failed to truncate nota_itens (run migrations first): %v", err)
	}
	if _, err := db.Exec("TRUNCATE TABLE notas RESTART IDENTITY CASCADE"); err != nil {
		t.Fatalf("failed to truncate notas (run migrations first): %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}
