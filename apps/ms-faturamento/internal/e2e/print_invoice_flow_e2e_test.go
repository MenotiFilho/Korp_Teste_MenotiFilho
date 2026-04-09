package e2e

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	stockDBURL   = "postgres://postgres:postgres@localhost:5433/estoque?sslmode=disable"
	billingDBURL = "postgres://postgres:postgres@localhost:5434/faturamento?sslmode=disable"
)

type testService struct {
	name   string
	cmd    *exec.Cmd
	logs   *bytes.Buffer
	waitCh chan error
}

type createInvoiceResponse struct {
	ID int64 `json:"id"`
}

type listInvoiceResponseItem struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

type listProductResponseItem struct {
	Codigo string `json:"codigo"`
	Saldo  int    `json:"saldo"`
}

func TestPrintInvoiceFlow_EndToEndAcrossServices_ShouldCloseInvoiceAndDecreaseStock(t *testing.T) {
	if testing.Short() {
		t.Skip("e2e test skipped in short mode")
	}

	// Arrange
	ctx := context.Background()
	stockPort := findFreePort(t)
	billingPort := findFreePort(t)
	stockBaseURL := fmt.Sprintf("http://127.0.0.1:%d", stockPort)
	billingBaseURL := fmt.Sprintf("http://127.0.0.1:%d", billingPort)

	if err := pingDB(stockDBURL); err != nil {
		t.Fatalf("failed to connect stock DB (%s): %v. Ensure docker DB is up and migrations are applied.", stockDBURL, err)
	}
	if err := pingDB(billingDBURL); err != nil {
		t.Fatalf("failed to connect billing DB (%s): %v. Ensure docker DB is up and migrations are applied.", billingDBURL, err)
	}

	resetStockDB(t, ctx)
	resetBillingDB(t, ctx)

	msFaturamentoDir, msEstoqueDir := resolveServiceDirs(t)
	stockSvc := startService(t, "ms-estoque", msEstoqueDir, map[string]string{
		"PORT":   fmt.Sprintf("%d", stockPort),
		"DB_URL": stockDBURL,
	}, stockBaseURL+"/health")
	billingSvc := startService(t, "ms-faturamento", msFaturamentoDir, map[string]string{
		"PORT":        fmt.Sprintf("%d", billingPort),
		"DB_URL":      billingDBURL,
		"ESTOQUE_URL": stockBaseURL,
	}, billingBaseURL+"/health")

	t.Cleanup(func() {
		stopService(t, billingSvc)
		stopService(t, stockSvc)
	})

	// Act 1: create stock product
	productPayload := `{"codigo":"P-E2E-001","descricao":"Produto E2E","saldo":10}`
	productRes, productBody := doJSONRequest(t, http.MethodPost, stockBaseURL+"/api/v1/produtos", productPayload, map[string]string{
		"X-Request-ID": "req-e2e-create-product",
	})
	if productRes.StatusCode != http.StatusCreated {
		t.Fatalf("expected create product status 201, got %d, body=%s", productRes.StatusCode, productBody)
	}

	// Act 2: create invoice in billing
	invoicePayload := `{"itens":[{"produto_codigo":"P-E2E-001","quantidade":3}]}`
	createInvRes, createInvBody := doJSONRequest(t, http.MethodPost, billingBaseURL+"/api/v1/notas", invoicePayload, map[string]string{
		"X-Request-ID": "req-e2e-create-invoice",
	})
	if createInvRes.StatusCode != http.StatusCreated {
		t.Fatalf("expected create invoice status 201, got %d, body=%s", createInvRes.StatusCode, createInvBody)
	}

	var createdInvoice createInvoiceResponse
	if err := json.Unmarshal([]byte(createInvBody), &createdInvoice); err != nil {
		t.Fatalf("failed to parse create invoice response: %v; body=%s", err, createInvBody)
	}
	if createdInvoice.ID <= 0 {
		t.Fatalf("expected created invoice id > 0, got %d", createdInvoice.ID)
	}

	// Act 3: print invoice (critical integration path)
	printURL := fmt.Sprintf("%s/api/v1/notas/%d/imprimir", billingBaseURL, createdInvoice.ID)
	printRes, printBody := doJSONRequest(t, http.MethodPost, printURL, "", map[string]string{
		"X-Request-ID": "req-e2e-print-invoice",
	})
	if printRes.StatusCode != http.StatusOK {
		t.Fatalf("expected print status 200, got %d, body=%s\nms-faturamento logs:\n%s\nms-estoque logs:\n%s", printRes.StatusCode, printBody, billingSvc.logs.String(), stockSvc.logs.String())
	}

	// Assert 1: invoice is closed
	listInvRes, listInvBody := doJSONRequest(t, http.MethodGet, billingBaseURL+"/api/v1/notas", "", map[string]string{
		"X-Request-ID": "req-e2e-list-invoices",
	})
	if listInvRes.StatusCode != http.StatusOK {
		t.Fatalf("expected list invoices status 200, got %d, body=%s", listInvRes.StatusCode, listInvBody)
	}

	var invoices []listInvoiceResponseItem
	if err := json.Unmarshal([]byte(listInvBody), &invoices); err != nil {
		t.Fatalf("failed to parse invoice list: %v; body=%s", err, listInvBody)
	}

	invoiceFound := false
	for _, inv := range invoices {
		if inv.ID == createdInvoice.ID {
			invoiceFound = true
			if inv.Status != "FECHADA" {
				t.Fatalf("expected invoice %d status FECHADA, got %s", createdInvoice.ID, inv.Status)
			}
		}
	}
	if !invoiceFound {
		t.Fatalf("expected invoice id %d to be listed", createdInvoice.ID)
	}

	// Assert 2: stock was decreased in stock service DB path
	listProdRes, listProdBody := doJSONRequest(t, http.MethodGet, stockBaseURL+"/api/v1/produtos", "", map[string]string{
		"X-Request-ID": "req-e2e-list-products",
	})
	if listProdRes.StatusCode != http.StatusOK {
		t.Fatalf("expected list products status 200, got %d, body=%s", listProdRes.StatusCode, listProdBody)
	}

	var products []listProductResponseItem
	if err := json.Unmarshal([]byte(listProdBody), &products); err != nil {
		t.Fatalf("failed to parse product list: %v; body=%s", err, listProdBody)
	}

	productFound := false
	for _, p := range products {
		if p.Codigo == "P-E2E-001" {
			productFound = true
			if p.Saldo != 7 {
				t.Fatalf("expected stock saldo 7 after print, got %d", p.Saldo)
			}
		}
	}
	if !productFound {
		t.Fatal("expected product P-E2E-001 to be listed")
	}
}

func resolveServiceDirs(t *testing.T) (msFaturamentoDir string, msEstoqueDir string) {
	t.Helper()

	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve test file path")
	}

	e2eDir := filepath.Dir(testFile)
	msFaturamentoDir = filepath.Clean(filepath.Join(e2eDir, "../.."))
	msEstoqueDir = filepath.Clean(filepath.Join(msFaturamentoDir, "../ms-estoque"))

	return msFaturamentoDir, msEstoqueDir
}

func startService(t *testing.T, name, dir string, overrides map[string]string, healthURL string) *testService {
	t.Helper()

	cmd := exec.Command("go", "run", "./cmd/"+name)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), mapToEnv(overrides)...)

	var logs bytes.Buffer
	cmd.Stdout = &logs
	cmd.Stderr = &logs

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start %s: %v", name, err)
	}

	svc := &testService{name: name, cmd: cmd, logs: &logs, waitCh: make(chan error, 1)}
	go func() {
		svc.waitCh <- cmd.Wait()
	}()

	waitForHealth(t, svc, healthURL, 30*time.Second)
	return svc
}

func stopService(t *testing.T, svc *testService) {
	t.Helper()

	if svc == nil || svc.cmd == nil || svc.cmd.Process == nil {
		return
	}

	_ = svc.cmd.Process.Signal(os.Interrupt)

	select {
	case <-time.After(5 * time.Second):
		_ = svc.cmd.Process.Kill()
		select {
		case <-svc.waitCh:
		case <-time.After(2 * time.Second):
		}
	case <-svc.waitCh:
	}
}

func waitForHealth(t *testing.T, svc *testService, healthURL string, timeout time.Duration) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		select {
		case err := <-svc.waitCh:
			t.Fatalf("%s exited before healthcheck: %v\nlogs:\n%s", svc.name, err, svc.logs.String())
		default:
		}

		resp, err := http.Get(healthURL)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}

		time.Sleep(200 * time.Millisecond)
	}

	t.Fatalf("timeout waiting for %s health at %s\nlogs:\n%s", svc.name, healthURL, svc.logs.String())
}

func doJSONRequest(t *testing.T, method, url, body string, headers map[string]string) (*http.Response, string) {
	t.Helper()

	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		t.Fatalf("failed to create request %s %s: %v", method, url, err)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed %s %s: %v", method, url, err)
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body for %s %s: %v", method, url, err)
	}

	return resp, string(bodyBytes)
}

func mapToEnv(m map[string]string) []string {
	if len(m) == 0 {
		return nil
	}

	out := make([]string, 0, len(m))
	for k, v := range m {
		out = append(out, k+"="+v)
	}
	return out
}

func findFreePort(t *testing.T) int {
	t.Helper()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to allocate free port: %v", err)
	}
	defer l.Close()

	addr, ok := l.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatalf("unexpected listener addr type: %T", l.Addr())
	}
	return addr.Port
}

func pingDB(dbURL string) error {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return db.PingContext(ctx)
}

func resetStockDB(t *testing.T, ctx context.Context) {
	t.Helper()

	db, err := sql.Open("pgx", stockDBURL)
	if err != nil {
		t.Fatalf("failed to open stock db: %v", err)
	}
	defer db.Close()

	if _, err := db.ExecContext(ctx, "TRUNCATE TABLE estoque_baixas RESTART IDENTITY"); err != nil {
		t.Fatalf("failed to truncate estoque_baixas (run stock migrations first): %v", err)
	}
	if _, err := db.ExecContext(ctx, "TRUNCATE TABLE produtos RESTART IDENTITY CASCADE"); err != nil {
		t.Fatalf("failed to truncate produtos (run stock migrations first): %v", err)
	}
}

func resetBillingDB(t *testing.T, ctx context.Context) {
	t.Helper()

	db, err := sql.Open("pgx", billingDBURL)
	if err != nil {
		t.Fatalf("failed to open billing db: %v", err)
	}
	defer db.Close()

	if _, err := db.ExecContext(ctx, "TRUNCATE TABLE nota_itens RESTART IDENTITY CASCADE"); err != nil {
		t.Fatalf("failed to truncate nota_itens (run billing migrations first): %v", err)
	}
	if _, err := db.ExecContext(ctx, "TRUNCATE TABLE notas RESTART IDENTITY CASCADE"); err != nil {
		t.Fatalf("failed to truncate notas (run billing migrations first): %v", err)
	}
}
