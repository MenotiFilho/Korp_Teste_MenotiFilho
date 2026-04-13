package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"syscall"
	"time"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/domain"
)

var (
	ErrStockProductNotFound   = errors.New("stock: product not found")
	ErrStockInsufficientStock = errors.New("stock: insufficient stock")
	ErrEstoqueUnavailable     = errors.New("estoque service unavailable")
)

type StockClient struct {
	baseURL    string
	httpClient *http.Client
	circuit    *CircuitBreaker
	retryDelay time.Duration
}

func NewStockClient(baseURL string, timeout time.Duration) *StockClient {
	return &StockClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		circuit:    NewCircuitBreaker(3, 10*time.Second),
		retryDelay: 500 * time.Millisecond,
	}
}

type decreaseStockRequest struct {
	Itens []decreaseStockItemRequest `json:"itens"`
}

type decreaseStockItemRequest struct {
	Codigo     string `json:"codigo"`
	Quantidade int    `json:"quantidade"`
}

type stockErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (c *StockClient) DecreaseStock(ctx context.Context, items []domain.StockDecreaseItem, idempotencyKey string) error {
	if err := c.circuit.Allow(); err != nil {
		return ErrEstoqueUnavailable
	}

	payload := decreaseStockRequest{
		Itens: make([]decreaseStockItemRequest, 0, len(items)),
	}
	for _, item := range items {
		payload.Itens = append(payload.Itens, decreaseStockItemRequest{
			Codigo:     item.Codigo,
			Quantidade: item.Quantidade,
		})
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal stock request: %w", err)
	}

	url := c.baseURL + "/api/v1/estoque/baixa"
	err = c.doRequest(ctx, url, body, idempotencyKey)
	if err != nil {
		if isTransientError(err) {
			select {
			case <-time.After(c.retryDelay):
			case <-ctx.Done():
				return ErrEstoqueUnavailable
			}
			err = c.doRequest(ctx, url, body, idempotencyKey)
		}
	}

	if err != nil {
		if isTransientError(err) {
			c.circuit.RecordFailure()
		}
		return err
	}

	c.circuit.RecordSuccess()
	return nil
}

func isTransientError(err error) bool {
	return errors.Is(err, ErrEstoqueUnavailable) &&
		!errors.Is(err, ErrStockProductNotFound) &&
		!errors.Is(err, ErrStockInsufficientStock)
}

func (c *StockClient) doRequest(ctx context.Context, url string, body []byte, idempotencyKey string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return ErrEstoqueUnavailable
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", idempotencyKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if isConnectionError(err) {
			return fmt.Errorf("%w: connection error", ErrEstoqueUnavailable)
		}
		return ErrEstoqueUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusNotFound {
		var stockErr stockErrorResponse
		if json.Unmarshal(respBody, &stockErr) == nil && stockErr.Code == "PRODUCT_NOT_FOUND" {
			return ErrStockProductNotFound
		}
		return ErrStockProductNotFound
	}

	if resp.StatusCode == http.StatusConflict {
		var stockErr stockErrorResponse
		if json.Unmarshal(respBody, &stockErr) == nil && stockErr.Code == "INSUFFICIENT_STOCK" {
			return ErrStockInsufficientStock
		}
		return ErrStockInsufficientStock
	}

	return ErrEstoqueUnavailable
}

func isConnectionError(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}
	if errors.Is(err, syscall.ECONNREFUSED) || errors.Is(err, syscall.ECONNRESET) {
		return true
	}
	return false
}
