package domain

import "testing"

func TestNewInvoice_WhenInputIsValid_ShouldReturnInvoiceWithStatusAberta(t *testing.T) {
	// Arrange
	items := []InvoiceItem{
		{ProdutoCodigo: "P-001", Quantidade: 2},
	}

	// Act
	invoice, err := NewInvoice(items)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if invoice.Status != StatusAberta {
		t.Fatalf("expected status %q, got %q", StatusAberta, invoice.Status)
	}
	if len(invoice.Itens) != 1 {
		t.Fatalf("expected 1 item, got %d", len(invoice.Itens))
	}
	if invoice.Itens[0].ProdutoCodigo != "P-001" || invoice.Itens[0].Quantidade != 2 {
		t.Fatalf("unexpected item: %+v", invoice.Itens[0])
	}
}

func TestNewInvoice_WhenItemsEmpty_ShouldReturnError(t *testing.T) {
	// Arrange
	items := []InvoiceItem{}

	// Act
	_, err := NewInvoice(items)

	// Assert
	if err == nil {
		t.Fatal("expected error for empty items")
	}
	if err != ErrInvoiceItemsRequired {
		t.Fatalf("expected ErrInvoiceItemsRequired, got %v", err)
	}
}

func TestNewInvoice_WhenItemsNil_ShouldReturnError(t *testing.T) {
	// Arrange
	var items []InvoiceItem

	// Act
	_, err := NewInvoice(items)

	// Assert
	if err == nil {
		t.Fatal("expected error for nil items")
	}
	if err != ErrInvoiceItemsRequired {
		t.Fatalf("expected ErrInvoiceItemsRequired, got %v", err)
	}
}

func TestNewInvoiceItem_WhenInputIsValid_ShouldReturnItem(t *testing.T) {
	// Arrange
	codigo := "P-001"
	quantidade := 5

	// Act
	item, err := NewInvoiceItem(codigo, quantidade)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if item.ProdutoCodigo != codigo {
		t.Fatalf("expected codigo %q, got %q", codigo, item.ProdutoCodigo)
	}
	if item.Quantidade != quantidade {
		t.Fatalf("expected quantidade %d, got %d", quantidade, item.Quantidade)
	}
}

func TestNewInvoiceItem_WhenCodigoIsEmpty_ShouldReturnError(t *testing.T) {
	// Arrange
	quantidade := 5

	// Act
	_, err := NewInvoiceItem("   ", quantidade)

	// Assert
	if err == nil {
		t.Fatal("expected error for empty codigo")
	}
	if err != ErrProdutoCodigoRequired {
		t.Fatalf("expected ErrProdutoCodigoRequired, got %v", err)
	}
}

func TestNewInvoiceItem_WhenQuantidadeIsZero_ShouldReturnError(t *testing.T) {
	// Arrange
	codigo := "P-001"

	// Act
	_, err := NewInvoiceItem(codigo, 0)

	// Assert
	if err == nil {
		t.Fatal("expected error for zero quantidade")
	}
	if err != ErrQuantidadeMustBePositive {
		t.Fatalf("expected ErrQuantidadeMustBePositive, got %v", err)
	}
}

func TestNewInvoiceItem_WhenQuantidadeIsNegative_ShouldReturnError(t *testing.T) {
	// Arrange
	codigo := "P-001"

	// Act
	_, err := NewInvoiceItem(codigo, -1)

	// Assert
	if err == nil {
		t.Fatal("expected error for negative quantidade")
	}
	if err != ErrQuantidadeMustBePositive {
		t.Fatalf("expected ErrQuantidadeMustBePositive, got %v", err)
	}
}

func TestInvoice_Close_WhenStatusIsAberta_ShouldChangeStatusToFechada(t *testing.T) {
	// Arrange
	items := []InvoiceItem{
		{ProdutoCodigo: "P-001", Quantidade: 2},
	}
	invoice, _ := NewInvoice(items)

	// Act
	err := invoice.Close()

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if invoice.Status != StatusFechada {
		t.Fatalf("expected status %q, got %q", StatusFechada, invoice.Status)
	}
}

func TestInvoice_Close_WhenStatusIsFechada_ShouldReturnError(t *testing.T) {
	// Arrange
	items := []InvoiceItem{
		{ProdutoCodigo: "P-001", Quantidade: 2},
	}
	invoice, _ := NewInvoice(items)
	_ = invoice.Close()

	// Act
	err := invoice.Close()

	// Assert
	if err == nil {
		t.Fatal("expected error closing already closed invoice")
	}
	if err != ErrInvoiceNotAberta {
		t.Fatalf("expected ErrInvoiceNotAberta, got %v", err)
	}
}

func TestInvoice_StockDecreaseItems_ShouldMapItemsToCodesAndQuantities(t *testing.T) {
	// Arrange
	items := []InvoiceItem{
		{ProdutoCodigo: "P-001", Quantidade: 2},
		{ProdutoCodigo: "P-002", Quantidade: 5},
	}
	invoice, _ := NewInvoice(items)

	// Act
	result := invoice.StockDecreaseItems()

	// Assert
	if len(result) != 2 {
		t.Fatalf("expected 2 items, got %d", len(result))
	}
	if result[0].Codigo != "P-001" || result[0].Quantidade != 2 {
		t.Fatalf("unexpected first item: %+v", result[0])
	}
	if result[1].Codigo != "P-002" || result[1].Quantidade != 5 {
		t.Fatalf("unexpected second item: %+v", result[1])
	}
}

func TestValidateInvoiceUpdate_WhenStatusIsAbertaAndItemsValid_ShouldReturnNil(t *testing.T) {
	// Arrange
	status := StatusAberta
	items := []InvoiceItem{
		{ProdutoCodigo: "P-001", Quantidade: 2},
	}

	// Act
	err := ValidateInvoiceUpdate(status, items)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateInvoiceUpdate_WhenStatusIsFechada_ShouldReturnError(t *testing.T) {
	// Arrange
	status := StatusFechada
	items := []InvoiceItem{
		{ProdutoCodigo: "P-001", Quantidade: 2},
	}

	// Act
	err := ValidateInvoiceUpdate(status, items)

	// Assert
	if err == nil {
		t.Fatal("expected error for FECHADA status")
	}
	if err != ErrInvoiceNotAberta {
		t.Fatalf("expected ErrInvoiceNotAberta, got %v", err)
	}
}

func TestValidateInvoiceUpdate_WhenItemsEmpty_ShouldReturnError(t *testing.T) {
	// Arrange
	status := StatusAberta
	items := []InvoiceItem{}

	// Act
	err := ValidateInvoiceUpdate(status, items)

	// Assert
	if err == nil {
		t.Fatal("expected error for empty items")
	}
	if err != ErrInvoiceItemsRequired {
		t.Fatalf("expected ErrInvoiceItemsRequired, got %v", err)
	}
}

func TestValidateInvoiceDelete_WhenStatusIsAberta_ShouldReturnNil(t *testing.T) {
	// Arrange
	status := StatusAberta

	// Act
	err := ValidateInvoiceDelete(status)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateInvoiceDelete_WhenStatusIsFechada_ShouldReturnError(t *testing.T) {
	// Arrange
	status := StatusFechada

	// Act
	err := ValidateInvoiceDelete(status)

	// Assert
	if err == nil {
		t.Fatal("expected error for FECHADA status")
	}
	if err != ErrInvoiceNotAberta {
		t.Fatalf("expected ErrInvoiceNotAberta, got %v", err)
	}
}
