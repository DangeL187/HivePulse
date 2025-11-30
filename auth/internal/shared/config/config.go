package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/DangeL187/erax"
	"github.com/joho/godotenv"
)

type Config struct {
	PostgresDSN           string
	CasbinModelConfigPath string

	DBConnectTimeout      time.Duration
	DeviceAccessTokenTTL  time.Duration
	DeviceRefreshTokenTTL time.Duration
	UserAccessTokenTTL    time.Duration
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("[env] no .env file found, skipping")
	}

	postgresDSN, err := loadPostgresDSN()
	if err != nil {
		return nil, erax.Wrap(err, "failed to load postgres dsn")
	}

	cfg := &Config{
		PostgresDSN:           postgresDSN,
		CasbinModelConfigPath: "casbin_model.conf",
		DBConnectTimeout:      1 * time.Minute,
		DeviceAccessTokenTTL:  10 * time.Minute,
		DeviceRefreshTokenTTL: 24 * time.Hour,
		UserAccessTokenTTL:    10 * time.Minute,
	}

	return cfg, nil
}

func loadPostgresDSN() (string, error) {
	vars := map[string]string{
		"POSTGRES_HOST":     os.Getenv("POSTGRES_HOST"),
		"POSTGRES_PORT":     os.Getenv("POSTGRES_PORT"),
		"POSTGRES_USER":     os.Getenv("POSTGRES_USER"),
		"POSTGRES_PASSWORD": os.Getenv("POSTGRES_PASSWORD"),
		"POSTGRES_DB":       os.Getenv("POSTGRES_DB"),
		"POSTGRES_SSL_MODE": os.Getenv("POSTGRES_SSL_MODE"),
	}

	for name, value := range vars {
		if value == "" {
			return "", fmt.Errorf("missing required env var: %s", name)
		}
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		vars["POSTGRES_HOST"],
		vars["POSTGRES_PORT"],
		vars["POSTGRES_USER"],
		vars["POSTGRES_PASSWORD"],
		vars["POSTGRES_DB"],
		vars["POSTGRES_SSL_MODE"],
	)

	return dsn, nil
}
