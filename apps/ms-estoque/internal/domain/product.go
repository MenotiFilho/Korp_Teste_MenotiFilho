package domain

import (
	"errors"
	"strings"
)

var (
	ErrCodigoRequired    = errors.New("codigo is required")
	ErrDescricaoRequired = errors.New("descricao is required")
	ErrSaldoNegative     = errors.New("saldo cannot be negative")
)

type Product struct {
	ID        int64
	Codigo    string
	Descricao string
	Saldo     int
}

func NewProduct(codigo, descricao string, saldo int) (Product, error) {
	codigo = strings.TrimSpace(codigo)
	descricao = strings.TrimSpace(descricao)

	if codigo == "" {
		return Product{}, ErrCodigoRequired
	}
	if descricao == "" {
		return Product{}, ErrDescricaoRequired
	}
	if saldo < 0 {
		return Product{}, ErrSaldoNegative
	}

	return Product{
		Codigo:    codigo,
		Descricao: descricao,
		Saldo:     saldo,
	}, nil
}
