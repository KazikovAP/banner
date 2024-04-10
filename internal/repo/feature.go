package repo

import (
	"banner/internal/lib/logger"
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

func NewFeature(db *pgxpool.Pool, log *slog.Logger) *FeatureRepo {
	return &FeatureRepo{db, log}
}

func (f *FeatureRepo) CreateFeature(ctx context.Context, feature *models.Feature) error {
	_, err := f.db.Exec(ctx, `INSERT INTO features (name) VALUES ($1)`, feature.Name)
	if err != nil {
		f.log.Error("Failed to create Feature", logger.Err(err))
	}

	return nil
}

func (f *FeatureRepo) FindFeatureName(ctx context.Context, name string) (models.Feature, error) {
	query, err := f.db.Query(ctx, `SELECT * FROM features WHERE name = $1`, name)
	if err != nil {
		f.log.Error("Feature not found", logger.Err(err))
		return models.Feature{}, err
	}

	order := models.Feature{}
	if !query.Next() {
		f.log.Error("Feature not found")
		return models.Feature{}, fmt.Errorf("Feature not found")
	} else {
		err := query.Scan(&order.ID, &order.Name)
		if err != nil {
			f.log.Error("Feature not found", logger.Err(err))
		}
	}

	return order, nil
}

func (f *FeatureRepo) FindFeatureId(ctx context.Context, id int) (models.Feature, error) {
	var order models.Feature
	err := f.db.QueryRow(ctx, `SELECT id, name FROM features WHERE id = $1`, id).Scan(&order.ID, &order.Name)
	if err != nil {
		f.log.Error("Failed to find Feature by ID", logger.Err(err))
		return models.Feature{}, err
	}

	return order, nil
}
