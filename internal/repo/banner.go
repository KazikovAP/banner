package repo

import (
	"banner/internal/lib/logger"
	"banner/internal/models"
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BannerRepo struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewBanner(db *pgxpool.Pool, log *slog.Logger) *BannerRepo {
	return &BannerRepo{db, log}
}

func (b *BannerRepo) CreateBanner(ctx context.Context, banner *models.Banner) error {
	_, err := b.db.Exec(ctx,
		`INSERT INTO banners (feature_id, content, is_active, created_at, updated_at) VALUES ($1,$2,$3,$4,$5)`,
		banner.FeatureID, banner.Content, banner.IsActive, banner.CreatedAt, banner.UpdatedAt)

	if err != nil {
		b.log.Error("Failed to create banner", logger.Err(err))
		return err
	}

	return nil
}

func (b *BannerRepo) FindBannerId(ctx context.Context, id int) (models.Banner, error) {
	query, err := b.db.Query(ctx, `SELECT * FROM banners WHERE id = $1`, id)
	if err != nil {
		b.log.Error("Banner not found", logger.Err(err))
		return models.Banner{}, err
	}
	defer query.Close()

	order := models.Banner{}
	for query.Next() {
		err := query.Scan(&order.ID, &order.FeatureID, &order.Content, &order.IsActive, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			b.log.Error("Banner not found", logger.Err(err))
			return models.Banner{}, err
		}
	}

	return order, nil
}

func (b *BannerRepo) FindBannerFeatureID(ctx context.Context, feature_id int) ([]models.Banner, error) {
	query, err := b.db.Query(ctx, `SELECT * FROM banners WHERE feature_id = $1`, feature_id)
	if err != nil {
		b.log.Error("Error querying banners", logger.Err(err))
		return nil, err
	}
	defer query.Close()

	var result []models.Banner
	for query.Next() {
		var order models.Banner
		err := query.Scan(&order.ID, &order.FeatureID, &order.Content, &order.IsActive, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			b.log.Error("Error scanning banners", logger.Err(err))
			return nil, err
		}
		result = append(result, order)
	}

	if len(result) == 0 {
		b.log.Info("No banners found for feature ID:", feature_id)
		return []models.Banner{}, nil
	}

	return result, nil
}

func (b *BannerRepo) FindBannerTagID(ctx context.Context, tagId int) ([]models.Banner, error) {
	query, err := b.db.Query(ctx, `SELECT * FROM banner_tags WHERE tag_id = $1`, tagId)
	if err != nil {
		b.log.Error("Error querying banners", logger.Err(err))
		return nil, err
	}
	defer query.Close()

	var result []models.Banner
	for query.Next() {
		var order models.Banner
		err := query.Scan(&order.ID, &order.FeatureID, &order.Content, &order.IsActive, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			b.log.Error("Error scanning banners", logger.Err(err))
			return nil, err
		}
		result = append(result, order)
	}

	if len(result) == 0 {
		b.log.Info("No banners found for tag ID:", tagId)
		return []models.Banner{}, nil
	}

	return result, nil
}
