package repo

import (
	logerr "banner/internal/lib/logger/logerr"
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
	err := b.db.QueryRow(ctx,
		`INSERT INTO banners (feature_id, content, is_active, created_at, updated_at) VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		banner.FeatureID, banner.Content, banner.IsActive, banner.CreatedAt, banner.UpdatedAt).Scan(&banner.ID)

	if err != nil {
		b.log.Error("Failed to create banner", logerr.Err(err))
		return err
	}

	return nil
}

func (b *BannerRepo) FindBannerId(ctx context.Context, id int) (models.Banner, error) {
	query, err := b.db.Query(ctx, `SELECT * FROM banners WHERE id = $1`, id)
	if err != nil {
		b.log.Error("Banner not found", logerr.Err(err))
		return models.Banner{}, err
	}
	defer query.Close()

	resultArray := models.Banner{}
	for query.Next() {
		err := query.Scan(&resultArray.ID, &resultArray.FeatureID, &resultArray.Content, &resultArray.IsActive, &resultArray.CreatedAt, &resultArray.UpdatedAt)
		if err != nil {
			b.log.Error("Banner not found", logerr.Err(err))
			return models.Banner{}, err
		}
	}

	return resultArray, nil
}

func (b *BannerRepo) FindBannersFeatureID(ctx context.Context, feature_id int) ([]models.Banner, error) {
	query, err := b.db.Query(ctx, `SELECT * FROM banners WHERE feature_id = $1`, feature_id)
	if err != nil {
		b.log.Error("Error querying banners", logerr.Err(err))
		return nil, err
	}
	defer query.Close()

	var resultSlice []models.Banner
	for query.Next() {
		var resultArray models.Banner
		err := query.Scan(&resultArray.ID, &resultArray.FeatureID, &resultArray.Content, &resultArray.IsActive, &resultArray.CreatedAt, &resultArray.UpdatedAt)
		if err != nil {
			b.log.Error("Error scanning banners", logerr.Err(err))
			return nil, err
		}
		resultSlice = append(resultSlice, resultArray)
	}

	if len(resultSlice) == 0 {
		b.log.Info("No banners found for feature ID:", feature_id)
		return []models.Banner{}, nil
	}

	return resultSlice, nil
}

func (b *BannerRepo) FindBannersTagID(ctx context.Context, tagId int) ([]models.Banner, error) {
	query, err := b.db.Query(ctx, `SELECT * FROM banner_tags WHERE tag_id = $1`, tagId)
	if err != nil {
		b.log.Error("Error querying banners", logerr.Err(err))
		return nil, err
	}
	defer query.Close()

	var resultSlice []models.Banner
	for query.Next() {
		var resultArray models.Banner
		err := query.Scan(&resultArray.ID, &resultArray.FeatureID, &resultArray.Content, &resultArray.IsActive, &resultArray.CreatedAt, &resultArray.UpdatedAt)
		if err != nil {
			b.log.Error("Error scanning banners", logerr.Err(err))
			return nil, err
		}
		resultSlice = append(resultSlice, resultArray)
	}

	if len(resultSlice) == 0 {
		b.log.Info("No banners found for tag ID:", tagId)
		return []models.Banner{}, nil
	}

	return resultSlice, nil
}
