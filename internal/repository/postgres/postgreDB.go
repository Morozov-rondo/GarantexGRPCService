package postgres

import (
	"fmt"

	"garantexGRPC/configs"
	"github.com/jmoiron/sqlx"
)

func NewPostgresDB(cfg configs.DatabaseConfig) (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password)

	db, err := sqlx.Open("postgres", dataSourceName)

	if err != nil {
		return nil, err
	}

	return db, nil
}
