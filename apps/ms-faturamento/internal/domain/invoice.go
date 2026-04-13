package domain

import (
	"errors"
	"strings"
)

const (
	StatusAberta  = "ABERTA"
	StatusFechada = "FECHADA"
)

var (
	ErrInvoiceItemsRequired     = errors.New("invoice items required")
	ErrProdutoCodigoRequired    = errors.New("produto codigo is required")
	ErrQuantidadeMustBePositive = errors.New("quantidade must be positive")
	ErrInvoiceNotAberta         = errors.New("invoice is not in ABERTA status")
)

type Invoice struct {
	ID     int64
	Numero int
	Status string
	Itens  []InvoiceItem
}

type InvoiceItem struct {
	ID            int64
	NotaID        int64
	ProdutoCodigo string
	Quantidade    int
}

func NewInvoice(items []InvoiceItem) (Invoice, error) {
	if len(items) == 0 {
		return Invoice{}, ErrInvoiceItemsRequired
	}

	normalized := make([]InvoiceItem, 0, len(items))
	for _, item := range items {
		ni, err := NewInvoiceItem(item.ProdutoCodigo, item.Quantidade)
		if err != nil {
			return Invoice{}, err
		}
		normalized = append(normalized, ni)
	}

	return Invoice{
		Status: StatusAberta,
		Itens:  normalized,
	}, nil
}

func NewInvoiceItem(codigo string, quantidade int) (InvoiceItem, error) {
	codigo = strings.TrimSpace(codigo)

	if codigo == "" {
		return InvoiceItem{}, ErrProdutoCodigoRequired
	}
	if quantidade <= 0 {
		return InvoiceItem{}, ErrQuantidadeMustBePositive
	}

	return InvoiceItem{
		ProdutoCodigo: codigo,
		Quantidade:    quantidade,
	}, nil
}

func (i *Invoice) Close() error {
	if i.Status != StatusAberta {
		return ErrInvoiceNotAberta
	}

	i.Status = StatusFechada
	return nil
}

type StockDecreaseItem struct {
	Codigo     string
	Quantidade int
}

func (i *Invoice) StockDecreaseItems() []StockDecreaseItem {
	result := make([]StockDecreaseItem, 0, len(i.Itens))
	for _, item := range i.Itens {
		result = append(result, StockDecreaseItem{
			Codigo:     item.ProdutoCodigo,
			Quantidade: item.Quantidade,
		})
	}
	return result
}

func ValidateInvoiceUpdate(status string, items []InvoiceItem) error {
	if status != StatusAberta {
		return ErrInvoiceNotAberta
	}
	if len(items) == 0 {
		return ErrInvoiceItemsRequired
	}
	for _, item := range items {
		if _, err := NewInvoiceItem(item.ProdutoCodigo, item.Quantidade); err != nil {
			return err
		}
	}
	return nil
}

func ValidateInvoiceDelete(status string) error {
	if status != StatusAberta {
		return ErrInvoiceNotAberta
	}
	return nil
}
