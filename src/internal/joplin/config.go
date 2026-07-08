package joplin

import (
	"os"
	"strconv"
)

// Config holds all runtime configuration for the Joplin client.
type Config struct {
	Token            string
	BaseURL          string
	TimeoutSeconds   int
	HTTPRetries      int
	HTTPRetryBackoff float64
}

// LoadConfig reads configuration from environment variables.
func LoadConfig() Config {
	return Config{
		Token:            os.Getenv("JOPLIN_TOKEN"),
		BaseURL:          envOr("JOPLIN_BASE_URL", "http://localhost:41184"),
		TimeoutSeconds:   envInt("JOPLIN_TIMEOUT_SECONDS", 30),
		HTTPRetries:      envInt("JOPLIN_HTTP_RETRIES", 3),
		HTTPRetryBackoff: envFloat("JOPLIN_HTTP_RETRY_BACKOFF_SECONDS", 1.0),
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func envFloat(key string, def float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return def
}
