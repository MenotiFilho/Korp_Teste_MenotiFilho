package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/domain"
)

func TestStockClient_DecreaseStock_WhenEstoqueReturns200_ShouldSucceed(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/estoque/baixa" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewStockClient(server.URL, 2*time.Second, 0)

	items := []domain.StockDecreaseItem{
		{Codigo: "P-001", Quantidade: 2},
	}

	// Act
	err := client.DecreaseStock(context.Background(), items)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestStockClient_DecreaseStock_WhenEstoqueReturns404_ShouldReturnProductNotFoundError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code":"PRODUCT_NOT_FOUND","message":"produto nao encontrado","request_id":"test"}`))
	}))
	defer server.Close()

	client := NewStockClient(server.URL, 2*time.Second, 0)

	items := []domain.StockDecreaseItem{
		{Codigo: "P-999", Quantidade: 1},
	}

	// Act
	err := client.DecreaseStock(context.Background(), items)

	// Assert
	if err == nil {
		t.Fatal("expected error for product not found")
	}
	if err != ErrStockProductNotFound {
		t.Fatalf("expected ErrStockProductNotFound, got %v", err)
	}
}

func TestStockClient_DecreaseStock_WhenEstoqueReturns409_ShouldReturnInsufficientStockError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"code":"INSUFFICIENT_STOCK","message":"saldo insuficiente","request_id":"test"}`))
	}))
	defer server.Close()

	client := NewStockClient(server.URL, 2*time.Second, 0)

	items := []domain.StockDecreaseItem{
		{Codigo: "P-001", Quantidade: 999},
	}

	// Act
	err := client.DecreaseStock(context.Background(), items)

	// Assert
	if err == nil {
		t.Fatal("expected error for insufficient stock")
	}
	if err != ErrStockInsufficientStock {
		t.Fatalf("expected ErrStockInsufficientStock, got %v", err)
	}
}

func TestStockClient_DecreaseStock_WhenEstoqueReturns500_ShouldReturnEstoqueUnavailableError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"code":"INTERNAL_ERROR","message":"erro interno","request_id":"test"}`))
	}))
	defer server.Close()

	client := NewStockClient(server.URL, 2*time.Second, 0)

	items := []domain.StockDecreaseItem{
		{Codigo: "P-001", Quantidade: 1},
	}

	// Act
	err := client.DecreaseStock(context.Background(), items)

	// Assert
	if err == nil {
		t.Fatal("expected error for internal server error")
	}
	if err != ErrEstoqueUnavailable {
		t.Fatalf("expected ErrEstoqueUnavailable, got %v", err)
	}
}

func TestStockClient_DecreaseStock_WhenEstoqueTimeout_ShouldReturnEstoqueUnavailableError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewStockClient(server.URL, 100*time.Millisecond, 0)

	items := []domain.StockDecreaseItem{
		{Codigo: "P-001", Quantidade: 1},
	}

	// Act
	err := client.DecreaseStock(context.Background(), items)

	// Assert
	if err == nil {
		t.Fatal("expected error for timeout")
	}
	if err != ErrEstoqueUnavailable {
		t.Fatalf("expected ErrEstoqueUnavailable, got %v", err)
	}
}

func TestStockClient_DecreaseStock_WhenEstoqueUnreachable_ShouldReturnEstoqueUnavailableError(t *testing.T) {
	// Arrange
	client := NewStockClient("http://192.0.2.1:9999", 500*time.Millisecond, 0)

	items := []domain.StockDecreaseItem{
		{Codigo: "P-001", Quantidade: 1},
	}

	// Act
	err := client.DecreaseStock(context.Background(), items)

	// Assert
	if err == nil {
		t.Fatal("expected error for unreachable service")
	}
	if err != ErrEstoqueUnavailable {
		t.Fatalf("expected ErrEstoqueUnavailable, got %v", err)
	}
}

func TestStockClient_DecreaseStock_When500WithRetry_ShouldRetryOnce(t *testing.T) {
	// Arrange
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"code":"INTERNAL_ERROR","message":"erro","request_id":"test"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewStockClient(server.URL, 2*time.Second, 1)

	items := []domain.StockDecreaseItem{
		{Codigo: "P-001", Quantidade: 1},
	}

	// Act
	err := client.DecreaseStock(context.Background(), items)

	// Assert
	if err != nil {
		t.Fatalf("expected no error after retry, got %v", err)
	}
	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
}

func TestStockClient_DecreaseStock_WhenSentCorrectPayload_ShouldSendItemsInBody(t *testing.T) {
	// Arrange
	var receivedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, 4096)
		n, _ := r.Body.Read(buf)
		receivedBody = string(buf[:n])

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/estoque/baixa" {
			t.Errorf("expected /api/v1/estoque/baixa, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewStockClient(server.URL, 2*time.Second, 0)

	items := []domain.StockDecreaseItem{
		{Codigo: "P-001", Quantidade: 2},
		{Codigo: "P-002", Quantidade: 5},
	}

	// Act
	err := client.DecreaseStock(context.Background(), items)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if receivedBody == "" {
		t.Fatal("expected body to be sent")
	}
	if !containsString(receivedBody, "P-001") {
		t.Fatalf("expected body to contain P-001, got %s", receivedBody)
	}
	if !containsString(receivedBody, "P-002") {
		t.Fatalf("expected body to contain P-002, got %s", receivedBody)
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || containsIndex(s, substr))
}

func containsIndex(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
