package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	defaultPort            = "8081"
	defaultReadTimeoutSec  = 10
	defaultWriteTimeoutSec = 10
	defaultIdleTimeoutSec  = 30
	defaultShutdownSec     = 10
)

type Config struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		Port:            getEnv("PORT", defaultPort),
		ReadTimeout:     secondsEnv("HTTP_READ_TIMEOUT_SEC", defaultReadTimeoutSec),
		WriteTimeout:    secondsEnv("HTTP_WRITE_TIMEOUT_SEC", defaultWriteTimeoutSec),
		IdleTimeout:     secondsEnv("HTTP_IDLE_TIMEOUT_SEC", defaultIdleTimeoutSec),
		ShutdownTimeout: secondsEnv("HTTP_SHUTDOWN_TIMEOUT_SEC", defaultShutdownSec),
	}

	if cfg.Port == "" {
		return Config{}, fmt.Errorf("PORT cannot be empty")
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
