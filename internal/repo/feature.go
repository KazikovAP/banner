package repo

import (
	logerr "banner/internal/lib/logger/logerr"
	"banner/internal/models"
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FeatureRepo struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewFeatureRepo(db *pgxpool.Pool, log *slog.Logger) *FeatureRepo {
	return &FeatureRepo{db, log}
}

func (f *FeatureRepo) CreateFeature(ctx context.Context, feature *models.Feature) error {
	err := f.db.QueryRow(ctx, `INSERT INTO features (name) VALUES ($1) RETURNING id`, feature.Name).Scan(&feature.ID)
	if err != nil {
		f.log.Error("Failed to create feature", logerr.Err(err))
		return err
	}

	return nil
}

func (f *FeatureRepo) FindFeatureId(ctx context.Context, id int) (models.Feature, error) {
	var res models.Feature
	err := f.db.QueryRow(ctx, `SELECT id, name FROM features WHERE id = $1`, id).Scan(&res.ID, &res.Name)
	if err != nil {
		f.log.Error("Failed to find Feature by ID", logerr.Err(err))
		return models.Feature{}, err
	}

	return res, nil
}

func (f *FeatureRepo) FindFeatureByName(ctx context.Context, name string) (models.Feature, error) {
	query, err := f.db.Query(ctx, `SELECT * FROM features WHERE name = $1`, name)
	if err != nil {
		f.log.Error("Feature not found", logerr.Err(err))
		return models.Feature{}, err
	}

	res := models.Feature{}
	if !query.Next() {
		f.log.Error("Feature not found")
		return models.Feature{}, fmt.Errorf("Feature not found")
	} else {
		err := query.Scan(&res.ID, &res.Name)
		if err != nil {
			f.log.Error("Feature not found", logerr.Err(err))
		}
	}

	return res, nil
}
