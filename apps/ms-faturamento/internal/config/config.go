package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	defaultPort            = "8082"
	defaultDatabaseURL     = "postgres://postgres:postgres@localhost:5434/faturamento?sslmode=disable"
	defaultEstoqueURL      = "http://localhost:8081"
	defaultReadHeaderSec   = 5
	defaultMaxHeaderBytes  = 1 << 20
	defaultMaxBodyBytes    = 1 << 20
	defaultReadTimeoutSec  = 10
	defaultWriteTimeoutSec = 10
	defaultIdleTimeoutSec  = 30
	defaultShutdownSec     = 10
)

type Config struct {
	Port              string
	DatabaseURL       string
	EstoqueURL        string
	ReadHeaderTimeout time.Duration
	MaxHeaderBytes    int
	MaxBodyBytes      int64
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		Port:              getEnv("PORT", defaultPort),
		DatabaseURL:       getEnv("DB_URL", defaultDatabaseURL),
		EstoqueURL:        getEnv("ESTOQUE_URL", defaultEstoqueURL),
		ReadHeaderTimeout: secondsEnv("HTTP_READ_HEADER_TIMEOUT_SEC", defaultReadHeaderSec),
		MaxHeaderBytes:    intEnv("HTTP_MAX_HEADER_BYTES", defaultMaxHeaderBytes),
		MaxBodyBytes:      int64(intEnv("HTTP_MAX_BODY_BYTES", defaultMaxBodyBytes)),
		ReadTimeout:       secondsEnv("HTTP_READ_TIMEOUT_SEC", defaultReadTimeoutSec),
		WriteTimeout:      secondsEnv("HTTP_WRITE_TIMEOUT_SEC", defaultWriteTimeoutSec),
		IdleTimeout:       secondsEnv("HTTP_IDLE_TIMEOUT_SEC", defaultIdleTimeoutSec),
		ShutdownTimeout:   secondsEnv("HTTP_SHUTDOWN_TIMEOUT_SEC", defaultShutdownSec),
	}

	if cfg.Port == "" {
		return Config{}, fmt.Errorf("PORT cannot be empty")
	}
	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DB_URL cannot be empty")
	}
	if cfg.EstoqueURL == "" {
		return Config{}, fmt.Errorf("ESTOQUE_URL cannot be empty")
	}
	if cfg.MaxHeaderBytes <= 0 {
		return Config{}, fmt.Errorf("HTTP_MAX_HEADER_BYTES must be positive")
	}
	if cfg.MaxBodyBytes <= 0 {
		return Config{}, fmt.Errorf("HTTP_MAX_BODY_BYTES must be positive")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func secondsEnv(key string, fallback int) time.Duration {
	v := getEnv(key, strconv.Itoa(fallback))
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		n = fallback
	}
	return time.Duration(n) * time.Second
}

func intEnv(key string, fallback int) int {
	v := getEnv(key, strconv.Itoa(fallback))
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}
