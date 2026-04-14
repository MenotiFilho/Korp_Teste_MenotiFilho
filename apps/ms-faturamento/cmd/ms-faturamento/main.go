package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/config"
	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/httpapi"
	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/middleware"
	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/repository"
	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/service"
	"github.com/rs/cors"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.PingContext(context.Background()); err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	mux := httpapi.NewMux()
	invoiceRepo := repository.NewInvoiceRepository(db)
	invoiceService := service.NewInvoiceService(invoiceRepo)
	invoiceHandler := httpapi.NewInvoiceHandler(invoiceService)
	stockClient := service.NewStockClient(cfg.EstoqueURL, 2*time.Second)
	printService := service.NewPrintInvoiceService(invoiceRepo, stockClient)
	printHandler := httpapi.NewPrintInvoiceHandler(invoiceRepo, printService)
	httpapi.RegisterInvoiceRoutes(mux, invoiceHandler)
	httpapi.RegisterPrintInvoiceRoutes(mux, printHandler)
	// Wrap mux with project middlewares first
	handler := middleware.RequestID(middleware.Recover(middleware.MaxBodyBytes(cfg.MaxBodyBytes, middleware.Logger(mux))))

	// Apply CORS using rs/cors so the frontend (localhost:4200) can call the APIs
	c := cors.New(cors.Options{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Idempotency-Key", "X-Request-ID"},
		AllowCredentials: false,
	})

	handler = c.Handler(handler)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		MaxHeaderBytes:    cfg.MaxHeaderBytes,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		slog.Info("ms-faturamento starting", "port", cfg.Port)
		errCh <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server stopped with error", "error", err)
			os.Exit(1)
		}
		return
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	slog.Info("ms-faturamento stopped")
}
