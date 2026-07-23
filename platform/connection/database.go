package connection

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(ctx context.Context) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	var db *gorm.DB
	var err error

	for i := 0; i < 30; i++ {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("database connection cancelled: %w", ctx.Err())
		default:
		}

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

		if err == nil {
			slog.Info("connected to database")

			if err := db.Use(otelgorm.NewPlugin()); err != nil {
				return nil, fmt.Errorf("failed to add otelgorm plugin: %v", err)
			}

			return db, nil
		}

		slog.Warn("failed to connect to database, retrying", "error", err, "attempt", i+1)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to database after retries: %v", err)
}
