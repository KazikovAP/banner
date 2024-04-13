package repo

import (
	logerr "banner/internal/lib/logger/logerr"
	"banner/internal/models"
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BannerTagRepo struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewBannerTagRepo(db *pgxpool.Pool, log *slog.Logger) *BannerTagRepo {
	return &BannerTagRepo{db, log}
}

func (bt *BannerTagRepo) CreateBannerTag(ctx context.Context, bannerTag *models.BannerTag) error {
	_, err := bt.db.Exec(ctx, `INSERT INTO banner_tags (banner_id, tag_id) VALUES ($1, $2)`, bannerTag.BannerID, bannerTag.TagID)
	if err != nil {
		bt.log.Error("Failed to create BannerTag", logerr.Err(err))
		return err
	}

	return nil
}

func (bt *BannerTagRepo) FindBannerTagBannerID(ctx context.Context, bannerID int) ([]models.BannerTag, error) {
	rows, err := bt.db.Query(ctx, `SELECT * FROM banner_tags WHERE banner_id = $1`, bannerID)
	if err != nil {
		bt.log.Error("Failed to find BannerTags by Banner ID", logerr.Err(err))
		return nil, err
	}
	defer rows.Close()

	var bannerTags []models.BannerTag
	for rows.Next() {
		var bannerTag models.BannerTag
		if err := rows.Scan(&bannerTag.BannerID, &bannerTag.TagID); err != nil {
			bt.log.Error("Failed to scan BannerTag row", logerr.Err(err))
			return nil, err
		}
		bannerTags = append(bannerTags, bannerTag)
	}

	if err := rows.Err(); err != nil {
		bt.log.Error("Error occurred while iterating BannerTag rows", logerr.Err(err))
		return nil, err
	}

	return bannerTags, nil
}
