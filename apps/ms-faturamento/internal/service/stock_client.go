package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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
	maxRetries int
}

func NewStockClient(baseURL string, timeout time.Duration, maxRetries int) *StockClient {
	return &StockClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		maxRetries: maxRetries,
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

func (c *StockClient) DecreaseStock(ctx context.Context, items []domain.StockDecreaseItem) error {
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
	maxAttempts := 1 + c.maxRetries

	var lastErr error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		lastErr = c.doRequest(ctx, url, body)
		if lastErr == nil {
			return nil
		}

		if !isTransientError(lastErr) {
			return lastErr
		}

		if attempt < maxAttempts-1 {
			time.Sleep(200 * time.Millisecond)
		}
	}

	return lastErr
}

func (c *StockClient) doRequest(ctx context.Context, url string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return ErrEstoqueUnavailable
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
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

func isTransientError(err error) bool {
	return errors.Is(err, ErrEstoqueUnavailable)
}
