package repository

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/domain"
)

const defaultStockTestDatabaseURL = "postgres://postgres:postgres@localhost:5433/estoque?sslmode=disable"

func TestProductRepository_DecreaseStock_WhenEnoughBalance_ShouldUpdateSaldo(t *testing.T) {
	// Arrange
	db := openStockTestDB(t)
	repo := NewProductRepository(db)
	ctx := context.Background()

	p, err := domain.NewProduct("P-001", "Produto 1", 10)
	if err != nil {
		t.Fatalf("unexpected domain error: %v", err)
	}
	created, err := repo.CreateProduct(ctx, p)
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}

	// Act
	err = repo.DecreaseStock(ctx, []domain.StockDecreaseItem{{Codigo: created.Codigo, Quantidade: 2}}, "idem-1")
	if err != nil {
		t.Fatalf("unexpected decrease error: %v", err)
	}

	products, err := repo.ListProducts(ctx)
	if err != nil {
		t.Fatalf("unexpected list error: %v", err)
	}

	// Assert
	if len(products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(products))
	}
	if products[0].Saldo != 8 {
		t.Fatalf("expected saldo 8, got %d", products[0].Saldo)
	}
}

func TestProductRepository_DecreaseStock_WhenInsufficientBalance_ShouldNotUpdateAndReturnError(t *testing.T) {
	// Arrange
	db := openStockTestDB(t)
	repo := NewProductRepository(db)
	ctx := context.Background()

	p, err := domain.NewProduct("P-001", "Produto 1", 1)
	if err != nil {
		t.Fatalf("unexpected domain error: %v", err)
	}
	created, err := repo.CreateProduct(ctx, p)
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}

	// Act
	err = repo.DecreaseStock(ctx, []domain.StockDecreaseItem{{Codigo: created.Codigo, Quantidade: 2}}, "idem-2")

	// Assert
	if !errors.Is(err, ErrProductInsufficientStock) {
		t.Fatalf("expected ErrProductInsufficientStock, got %v", err)
	}

	products, listErr := repo.ListProducts(ctx)
	if listErr != nil {
		t.Fatalf("unexpected list error: %v", listErr)
	}
	if len(products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(products))
	}
	if products[0].Saldo != 1 {
		t.Fatalf("expected saldo unchanged at 1, got %d", products[0].Saldo)
	}
}

func TestProductRepository_DecreaseStock_WhenSameIdempotencyKeyReplayed_ShouldBeNoOp(t *testing.T) {
	// Arrange
	db := openStockTestDB(t)
	repo := NewProductRepository(db)
	ctx := context.Background()

	p, err := domain.NewProduct("P-001", "Produto 1", 10)
	if err != nil {
		t.Fatalf("unexpected domain error: %v", err)
	}
	created, err := repo.CreateProduct(ctx, p)
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}

	idemKey := "idem-replay-1"
	items := []domain.StockDecreaseItem{{Codigo: created.Codigo, Quantidade: 2}}

	// Act
	err = repo.DecreaseStock(ctx, items, idemKey)
	if err != nil {
		t.Fatalf("unexpected first decrease error: %v", err)
	}

	err = repo.DecreaseStock(ctx, items, idemKey)
	if err != nil {
		t.Fatalf("unexpected replay decrease error: %v", err)
	}

	products, err := repo.ListProducts(ctx)
	if err != nil {
		t.Fatalf("unexpected list error: %v", err)
	}

	// Assert
	if len(products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(products))
	}
	if products[0].Saldo != 8 {
		t.Fatalf("expected saldo 8 after replay no-op, got %d", products[0].Saldo)
	}

	var count int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM estoque_baixas WHERE idempotency_key = $1", idemKey).Scan(&count); err != nil {
		t.Fatalf("failed to count idempotency records: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected exactly 1 idempotency record, got %d", count)
	}
}

func openStockTestDB(t *testing.T) *sql.DB {
	t.Helper()

	if testing.Short() {
		t.Skip("integration test skipped in short mode")
	}

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = defaultStockTestDatabaseURL
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping database (%s): %v", dbURL, err)
	}

	if _, err := db.Exec("TRUNCATE TABLE produtos RESTART IDENTITY"); err != nil {
		t.Fatalf("failed to truncate produtos (run migrations first): %v", err)
	}
	if _, err := db.Exec("TRUNCATE TABLE estoque_baixas RESTART IDENTITY"); err != nil {
		t.Fatalf("failed to truncate estoque_baixas (run migrations first): %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}
