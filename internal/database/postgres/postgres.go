package postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testTask/internal/config"
)

type Database struct {
	db *gorm.DB
}

func New(dbConfig config.DatabaseConfig) (*Database, error) {
	const op = "database/postgres/postgres/New"

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: failed to open database: %w", op, err)
	}

	return &Database{db: db}, nil
}
