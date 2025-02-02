package db

import (
	"database/sql"
	"fmt"
	"music/pkg/config"

	_ "github.com/lib/pq"
)

type DB struct {
	PostgreSQL *sql.DB
}

func NewDB(cfg *config.Config) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{PostgreSQL: db}, nil
}

func (db *DB) Close() error {
	if db.PostgreSQL != nil {
		return db.PostgreSQL.Close()
	}
	return nil
}
