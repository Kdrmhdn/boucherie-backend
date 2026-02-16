package configs

import (
	"os"
	"strconv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Port   int
	DBPath string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	port := 8080
	if v := os.Getenv("PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			port = p
		}
	}

	dbPath := "boucherie.db"
	if v := os.Getenv("DB_PATH"); v != "" {
		dbPath = v
	}

	return &Config{
		Port:   port,
		DBPath: dbPath,
	}
}
