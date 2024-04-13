package repo

import (
	logerr "banner/internal/lib/logger/logerr"
	"banner/internal/models"
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TagRepo struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewTagRepo(db *pgxpool.Pool, log *slog.Logger) *TagRepo {
	return &TagRepo{db, log}
}

func (t *TagRepo) CreateTag(ctx context.Context, tag *models.Tag) error {
	err := t.db.QueryRow(ctx, `INSERT INTO tags (name) VALUES ($1) RETURNING id`, tag.Name).Scan(&tag.ID)
	if err != nil {
		t.log.Error("Failed to create tag", logerr.Err(err))
		return err
	}

	return nil
}

func (t *TagRepo) FindTagId(ctx context.Context, id int) (models.Tag, error) {
	var tag models.Tag
	err := t.db.QueryRow(ctx, `SELECT id, name FROM tags WHERE id = $1`, id).Scan(&tag.ID, &tag.Name)
	if err != nil {
		t.log.Error("Failed to find Tag by ID", logerr.Err(err))
		return models.Tag{}, err
	}

	return tag, nil
}

func (t *TagRepo) FindTagName(ctx context.Context, name string) (models.Tag, error) {
	query, err := t.db.Query(ctx, `SELECT * FROM tags WHERE name = $1`, name)
	if err != nil {
		t.log.Error("Tag not found", logerr.Err(err))
		return models.Tag{}, err
	}

	res := models.Tag{}
	if !query.Next() {
		t.log.Error("Tag not found")
		return models.Tag{}, err
	} else {
		err := query.Scan(&res.ID, &res.Name)
		if err != nil {
			t.log.Error("Tag not found", logerr.Err(err))
			return models.Tag{}, fmt.Errorf("Tag not found")
		}
	}

	return res, nil
}
