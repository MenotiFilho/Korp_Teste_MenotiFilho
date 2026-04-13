package domain

import "testing"

func TestNewProduct_WhenInputIsValid_ShouldReturnProduct(t *testing.T) {
	// Arrange
	codigo := "P-001"
	descricao := "Produto de teste"
	saldo := 10

	// Act
	product, err := NewProduct(codigo, descricao, saldo)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if product.Codigo != codigo {
		t.Fatalf("expected codigo %q, got %q", codigo, product.Codigo)
	}
	if product.Descricao != descricao {
		t.Fatalf("expected descricao %q, got %q", descricao, product.Descricao)
	}
	if product.Saldo != saldo {
		t.Fatalf("expected saldo %d, got %d", saldo, product.Saldo)
	}
}

func TestNewProduct_WhenCodigoIsEmpty_ShouldReturnError(t *testing.T) {
	// Arrange
	descricao := "Produto de teste"
	saldo := 10

	// Act
	_, err := NewProduct("   ", descricao, saldo)

	// Assert
	if err == nil {
		t.Fatal("expected error for empty codigo")
	}
	if err != ErrCodigoRequired {
		t.Fatalf("expected ErrCodigoRequired, got %v", err)
	}
}

func TestNewProduct_WhenDescricaoIsEmpty_ShouldReturnError(t *testing.T) {
	// Arrange
	codigo := "P-001"
	saldo := 10

	// Act
	_, err := NewProduct(codigo, "  ", saldo)

	// Assert
	if err == nil {
		t.Fatal("expected error for empty descricao")
	}
	if err != ErrDescricaoRequired {
		t.Fatalf("expected ErrDescricaoRequired, got %v", err)
	}
}

func TestNewProduct_WhenSaldoIsNegative_ShouldReturnError(t *testing.T) {
	// Arrange
	codigo := "P-001"
	descricao := "Produto de teste"

	// Act
	_, err := NewProduct(codigo, descricao, -1)

	// Assert
	if err == nil {
		t.Fatal("expected error for negative saldo")
	}
	if err != ErrSaldoNegative {
		t.Fatalf("expected ErrSaldoNegative, got %v", err)
	}
}

func TestValidateProductUpdate_WhenInputIsValid_ShouldReturnNil(t *testing.T) {
	// Arrange
	descricao := "Produto atualizado"
	saldo := 20

	// Act
	err := ValidateProductUpdate(descricao, saldo)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateProductUpdate_WhenDescricaoIsEmpty_ShouldReturnError(t *testing.T) {
	// Arrange
	saldo := 20

	// Act
	err := ValidateProductUpdate("  ", saldo)

	// Assert
	if err == nil {
		t.Fatal("expected error for empty descricao")
	}
	if err != ErrDescricaoRequired {
		t.Fatalf("expected ErrDescricaoRequired, got %v", err)
	}
}

func TestValidateProductUpdate_WhenSaldoIsNegative_ShouldReturnError(t *testing.T) {
	// Arrange
	descricao := "Produto atualizado"

	// Act
	err := ValidateProductUpdate(descricao, -1)

	// Assert
	if err == nil {
		t.Fatal("expected error for negative saldo")
	}
	if err != ErrSaldoNegative {
		t.Fatalf("expected ErrSaldoNegative, got %v", err)
	}
}
