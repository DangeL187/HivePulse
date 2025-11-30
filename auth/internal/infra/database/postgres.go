package database

import (
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"

	"github.com/DangeL187/erax"

	"auth/internal/shared/config"
)

func NewPostgres(cfg *config.Config) (*gorm.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.DBConnectTimeout)
	defer cancel()

	var db *gorm.DB
	var err error

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, erax.Wrap(err, "failed to connect to database within timeout")
		default:
			db, err = gorm.Open(postgres.Open(cfg.PostgresDSN), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			})
			if err == nil {
				return db, nil
			}
			log.Printf("[DB] Postgres not ready yet. Retrying in 2s...")
			<-ticker.C
		}
	}
}
