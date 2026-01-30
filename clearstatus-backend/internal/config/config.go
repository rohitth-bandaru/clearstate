package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port           string
	DBDriver       string
	DBDSN          string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	GoogleClientID string
}

// Load reads all configuration from environment variables. No hardcoded defaults.
// DB: use DB_DSN (single URL) or DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME (build URL from parts).
// Required: PORT, DB_DRIVER, JWT_SECRET; and either DB_DSN or (DB_HOST, DB_USER, DB_PASSWORD, DB_NAME).
// Optional: GOOGLE_CLIENT_ID. DB_PORT defaults to 3306 when using parts.
func Load() *Config {
	cfg := &Config{
		Port:           os.Getenv("PORT"),
		DBDriver:       os.Getenv("DB_DRIVER"),
		DBDSN:          os.Getenv("DB_DSN"),
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         os.Getenv("DB_PORT"),
		DBUser:         os.Getenv("DB_USER"),
		DBPassword:     os.Getenv("DB_PASSWORD"),
		DBName:         os.Getenv("DB_NAME"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		GoogleClientID: os.Getenv("GOOGLE_CLIENT_ID"),
	}
	if cfg.DBDSN == "" && cfg.DBHost != "" {
		if cfg.DBPort == "" {
			cfg.DBPort = "3306"
		}
		cfg.DBDSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	}
	return cfg
}

// Validate returns an error if required config is missing.
func (c *Config) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("PORT is required")
	}
	if c.DBDriver == "" {
		return fmt.Errorf("DB_DRIVER is required")
	}
	if c.DBDSN == "" {
		if c.DBHost == "" || c.DBUser == "" || c.DBPassword == "" || c.DBName == "" {
			return fmt.Errorf("either DB_DSN or all of DB_HOST, DB_USER, DB_PASSWORD, DB_NAME are required")
		}
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	return nil
}

// ConnectionInfoForLog returns a safe string for logging (no password). Use when DB is built from parts; otherwise returns "***".
func (c *Config) ConnectionInfoForLog() string {
	if c.DBHost != "" {
		port := c.DBPort
		if port == "" {
			port = "3306"
		}
		return fmt.Sprintf("%s@%s:%s/%s", c.DBUser, c.DBHost, port, c.DBName)
	}
	return "***"
}

func (c *Config) IsProd() bool {
	v, _ := strconv.ParseBool(os.Getenv("PRODUCTION"))
	return v
}
