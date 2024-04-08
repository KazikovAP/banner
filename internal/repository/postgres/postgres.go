package postgres

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	db *sqlx.DB
}

type BannerActions interface {
	GetUserBanner() (struct{}, error)
}

func NewPostgres() (*Postgres, error) {
	db, err := sqlx.Open("postgres", "host=localhost port=5436 user=user dbname=postgres password=password sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("database connection success")

	return &Postgres{
		db: db,
	}, nil
}
