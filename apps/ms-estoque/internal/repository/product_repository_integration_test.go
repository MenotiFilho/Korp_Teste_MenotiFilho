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

const defaultTestDatabaseURL = "postgres://postgres:postgres@localhost:5433/estoque?sslmode=disable"

func TestProductRepository_CreateAndListProducts_ShouldPersistAndReturnProducts(t *testing.T) {
	// Arrange
	db := openTestDB(t)
	repo := NewProductRepository(db)
	ctx := context.Background()

	p1, err := domain.NewProduct("P-001", "Produto 1", 10)
	if err != nil {
		t.Fatalf("unexpected error creating domain product: %v", err)
	}
	p2, err := domain.NewProduct("P-002", "Produto 2", 5)
	if err != nil {
		t.Fatalf("unexpected error creating domain product: %v", err)
	}

	// Act
	created1, err := repo.CreateProduct(ctx, p1)
	if err != nil {
		t.Fatalf("unexpected error creating product 1: %v", err)
	}
	created2, err := repo.CreateProduct(ctx, p2)
	if err != nil {
		t.Fatalf("unexpected error creating product 2: %v", err)
	}
	products, err := repo.ListProducts(ctx)
	if err != nil {
		t.Fatalf("unexpected error listing products: %v", err)
	}

	// Assert
	if created1.ID <= 0 || created2.ID <= 0 {
		t.Fatalf("expected generated ids, got %d and %d", created1.ID, created2.ID)
	}
	if len(products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(products))
	}
	if products[0].Codigo != "P-001" || products[1].Codigo != "P-002" {
		t.Fatalf("expected products ordered by insert with codigos P-001 and P-002, got %q and %q", products[0].Codigo, products[1].Codigo)
	}
}

func TestProductRepository_CreateProduct_WhenCodigoDuplicate_ShouldReturnSpecificError(t *testing.T) {
	// Arrange
	db := openTestDB(t)
	repo := NewProductRepository(db)
	ctx := context.Background()

	p1, err := domain.NewProduct("P-001", "Produto 1", 10)
	if err != nil {
		t.Fatalf("unexpected error creating domain product: %v", err)
	}
	p2, err := domain.NewProduct("P-001", "Produto duplicado", 2)
	if err != nil {
		t.Fatalf("unexpected error creating domain product: %v", err)
	}

	if _, err := repo.CreateProduct(ctx, p1); err != nil {
		t.Fatalf("unexpected error creating first product: %v", err)
	}

	// Act
	_, err = repo.CreateProduct(ctx, p2)

	// Assert
	if err == nil {
		t.Fatal("expected duplicate codigo error")
	}
	if !errors.Is(err, ErrProductCodigoAlreadyExists) {
		t.Fatalf("expected ErrProductCodigoAlreadyExists, got %v", err)
	}
}

func TestProductRepository_UpdateProduct_ShouldUpdateDescricaoAndSaldo(t *testing.T) {
	// Arrange
	db := openTestDB(t)
	repo := NewProductRepository(db)
	ctx := context.Background()

	p, err := domain.NewProduct("P-001", "Produto 1", 10)
	if err != nil {
		t.Fatalf("unexpected error creating domain product: %v", err)
	}
	created, err := repo.CreateProduct(ctx, p)
	if err != nil {
		t.Fatalf("unexpected error creating product: %v", err)
	}

	// Act
	updated, err := repo.UpdateProduct(ctx, created.ID, "Produto 1 Atualizado", 25)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error updating product: %v", err)
	}
	if updated.Descricao != "Produto 1 Atualizado" {
		t.Fatalf("expected descricao 'Produto 1 Atualizado', got %q", updated.Descricao)
	}
	if updated.Saldo != 25 {
		t.Fatalf("expected saldo 25, got %d", updated.Saldo)
	}
	if updated.Codigo != "P-001" {
		t.Fatalf("expected codigo 'P-001', got %q", updated.Codigo)
	}
}

func TestProductRepository_UpdateProduct_WhenProductNotFound_ShouldReturnError(t *testing.T) {
	// Arrange
	db := openTestDB(t)
	repo := NewProductRepository(db)
	ctx := context.Background()

	// Act
	_, err := repo.UpdateProduct(ctx, 9999, "Descricao", 10)

	// Assert
	if err == nil {
		t.Fatal("expected error for non-existent product")
	}
	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

func TestProductRepository_SoftDeleteProduct_ShouldSetDeletedAt(t *testing.T) {
	// Arrange
	db := openTestDB(t)
	repo := NewProductRepository(db)
	ctx := context.Background()

	p, err := domain.NewProduct("P-001", "Produto 1", 10)
	if err != nil {
		t.Fatalf("unexpected error creating domain product: %v", err)
	}
	created, err := repo.CreateProduct(ctx, p)
	if err != nil {
		t.Fatalf("unexpected error creating product: %v", err)
	}

	// Act
	err = repo.SoftDeleteProduct(ctx, created.ID)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error soft deleting product: %v", err)
	}
	products, err := repo.ListProducts(ctx)
	if err != nil {
		t.Fatalf("unexpected error listing products: %v", err)
	}
	if len(products) != 0 {
		t.Fatalf("expected 0 products after soft delete, got %d", len(products))
	}
}

func TestProductRepository_SoftDeleteProduct_WhenProductNotFound_ShouldReturnError(t *testing.T) {
	// Arrange
	db := openTestDB(t)
	repo := NewProductRepository(db)
	ctx := context.Background()

	// Act
	err := repo.SoftDeleteProduct(ctx, 9999)

	// Assert
	if err == nil {
		t.Fatal("expected error for non-existent product")
	}
	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
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

	if _, err := db.Exec("TRUNCATE TABLE produtos RESTART IDENTITY"); err != nil {
		t.Fatalf("failed to truncate produtos (run migrations first): %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}
