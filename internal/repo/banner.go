package repo

import (
	logerr "banner/internal/lib/logger/logerr"
	"banner/internal/models"
	"banner/internal/server/handlers/banners"
	"context"
	"errors"
	"log/slog"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BannerRepo struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewBannerRepo(db *pgxpool.Pool, log *slog.Logger) *BannerRepo {
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

func (b *BannerRepo) FindBannerFeatureTag(ctx context.Context, featureID, tagID int) (*models.Banner, error) {
	query := `SELECT b.id, b.feature_id, b.content, b.is_active, b.created_at, b.updated_at
			  FROM banners b
			  INNER JOIN banner_tags bt ON b.id = bt.banner_id
			  WHERE b.feature_id = $1 AND bt.tag_id = $2`

	row := b.db.QueryRow(ctx, query, featureID, tagID)

	var banner models.Banner

	err := row.Scan(&banner.ID, &banner.FeatureID, &banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt)
	if err != nil {
		b.log.Error("Error with database", logerr.Err(err))
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}

		b.log.Error("Failed to find banner", logerr.Err(err))
		return nil, err
	}

	return &banner, nil
}

func (b *BannerRepo) FindBannersParameters(ctx context.Context, params banners.RequestGetBanners) ([]models.Banner, error) {
	query := "SELECT b.id, b.feature_id, b.content, b.is_active, b.created_at, b.updated_at, array_agg(bt.tag_id) AS tag_ids FROM banners b LEFT JOIN banner_tags bt ON b.id = bt.banner_id WHERE 1=1"
	args := []interface{}{}

	if params.FeatureID != nil {
		query += " AND b.feature_id = $" + strconv.Itoa(len(args)+1)
		args = append(args, *params.FeatureID)
	}

	if params.TagID != nil {
		query += " AND b.id IN (SELECT banner_id FROM banner_tags WHERE tag_id = $" + strconv.Itoa(len(args)+1) + ")"
		args = append(args, *params.TagID)
	}

	query += " GROUP BY b.id"

	if params.Limit != nil {
		query += " LIMIT $" + strconv.Itoa(len(args)+1)
		args = append(args, *params.Limit)
	}

	if params.Offset != nil {
		query += " OFFSET $" + strconv.Itoa(len(args)+1)
		args = append(args, *params.Offset)
	}

	rows, err := b.db.Query(ctx, query, args...)
	if err != nil {
		b.log.Error("Failed to query banners", logerr.Err(err))
		return nil, err
	}
	defer rows.Close()

	var banners []models.Banner
	for rows.Next() {
		var banner models.Banner
		var tagIDs []int
		if err := rows.Scan(&banner.ID, &banner.FeatureID, &banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt, &tagIDs); err != nil {
			b.log.Error("Failed to scan banner row", logerr.Err(err))
			return nil, err
		}
		banner.TagIDs = tagIDs
		banners = append(banners, banner)
	}

	if err := rows.Err(); err != nil {
		b.log.Error("Error occurred while iterating banner rows", logerr.Err(err))
		return nil, err
	}

	return banners, nil
}

func (b *BannerRepo) UpdateBanner(ctx context.Context, banner *models.Banner) error {
	tx, err := b.db.Begin(ctx)
	if err != nil {
		b.log.Error("Failed to begin transaction", logerr.Err(err))
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `DELETE FROM banner_tags WHERE banner_id = $1`, banner.ID)
	if err != nil {
		b.log.Error("Failed to delete old tags for banner", logerr.Err(err))
		return err
	}

	for _, tagID := range banner.TagIDs {
		_, err = tx.Exec(ctx, `INSERT INTO banner_tags (banner_id, tag_id) VALUES ($1, $2)`, banner.ID, tagID)
		if err != nil {
			b.log.Error("Failed to insert tag for banner", logerr.Err(err))
			return err
		}
	}

	_, err = tx.Exec(ctx,
		`UPDATE banners SET feature_id = $1, content = $2, is_active = $3, updated_at = $4 WHERE id = $5`,
		banner.FeatureID, banner.Content, banner.IsActive, banner.UpdatedAt, banner.ID)
	if err != nil {
		b.log.Error("Failed to update banner", logerr.Err(err))
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		b.log.Error("Failed to commit transaction", logerr.Err(err))
		return err
	}

	return nil
}

func (b *BannerRepo) DeleteBannerID(ctx context.Context, id int) error {
	_, err := b.db.Exec(ctx, `DELETE FROM banners WHERE id = $1`, id)
	if err != nil {
		b.log.Error("Failed to delete banner by ID", logerr.Err(err))
		return err
	}

	return nil
}
