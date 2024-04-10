package repo

import (
	"banner/internal/lib/logger"
	"banner/internal/models"
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BannerTagRepo struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewBannerTag(db *pgxpool.Pool, log *slog.Logger) *BannerTagRepo {
	return &BannerTagRepo{db, log}
}

func (bt *BannerTagRepo) CreateBannerTag(ctx context.Context, bannerTag *models.BannerTag) error {
	_, err := bt.db.Exec(ctx, `INSERT INTO banner_tags (banner_id, tag_id) VALUES ($1, $2)`, bannerTag.BannerID, bannerTag.TagID)
	if err != nil {
		bt.log.Error("Failed to create BannerTag", logger.Err(err))
		return err
	}

	return nil
}

func (bt *BannerTagRepo) FindBannerTagsBannerID(ctx context.Context, bannerID int) ([]models.BannerTag, error) {
	order, err := bt.db.Query(ctx, `SELECT * FROM banner_tags WHERE banner_id = $1`, bannerID)
	if err != nil {
		bt.log.Error("Failed to find BannerTags by Banner ID", logger.Err(err))
		return nil, err
	}
	defer order.Close()

	var result []models.BannerTag
	for order.Next() {
		var res models.BannerTag
		if err := order.Scan(&res.BannerID, &res.TagID); err != nil {
			bt.log.Error("Failed to scan BannerTag", logger.Err(err))
			return nil, err
		}
		result = append(result, res)
	}

	if err := order.Err(); err != nil {
		bt.log.Error("Error occurred while iterating BannerTag", logger.Err(err))
		return nil, err
	}

	return result, nil
}
