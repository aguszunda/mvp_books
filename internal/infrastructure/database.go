package infrastructure

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		dbPort,
		os.Getenv("DB_NAME"),
	)

	var db *gorm.DB
	var err error

	for i := 0; i < 30; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

		if err == nil {
			log.Println("Connected to Database!")

			// Use OTEL GORM Middleware
			if err := db.Use(otelgorm.NewPlugin()); err != nil {
				return nil, fmt.Errorf("failed to add otelgorm plugin: %v", err)
			}

			return db, nil
		}

		log.Printf("Failed to connect to GORM DB: %v, retrying...", err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to database after retries: %v", err)
}
