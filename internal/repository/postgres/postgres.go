package postgres

import (
	logerr "banner/internal/lib/logger/logerr"
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	DB  *pgxpool.Pool
	log *slog.Logger
}

var pgInit *Postgres

func NewPostgres(ctx context.Context, cont string, log *slog.Logger) (*Postgres, error) {
	var pgOnce sync.Once

	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, cont)
		if err != nil {
			fmt.Errorf("Cannot to create connection pool", logerr.Err(err))
			return
		}

		pgInit = &Postgres{db, log}
		if err := CreateTable(ctx, db, log); err != nil {
			return
		}
	})

	return pgInit, nil
}

func CreateTable(ctx context.Context, db *pgxpool.Pool, log *slog.Logger) error {
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS banners (
			id SERIAL PRIMARY KEY,
			feature_id INTEGER,
			content JSONB,
			is_active BOOLEAN,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("Failed to create banners table", logerr.Err(err))
	}

	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS tags (
			id SERIAL PRIMARY KEY,
			name TEXT
		)
	`)
	if err != nil {
		return fmt.Errorf("Failed to create tags table", logerr.Err(err))
	}

	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS banner_tags (
			banner_id INTEGER,
			tag_id INTEGER,
			PRIMARY KEY (banner_id, tag_id),
			FOREIGN KEY (banner_id) REFERENCES banners(id) ON DELETE CASCADE,
			FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("Failed to create banner_tags table", logerr.Err(err))
	}

	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS features (
			id SERIAL PRIMARY KEY,
			name TEXT
		)
	`)
	if err != nil {
		return fmt.Errorf("Failed to create features table", logerr.Err(err))
	}

	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
		    id SERIAL PRIMARY KEY ,
			username TEXT,
			password TEXT,
			role TEXT
		)
	`)
	if err != nil {
		return fmt.Errorf("Failed to create users table", logerr.Err(err))
	}

	return nil
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.DB.Ping(ctx)
}

func (pg *Postgres) Close() {
	pg.DB.Close()
}
