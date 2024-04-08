package handler

import (
	"context"

	"banner/internal/models"
	"banner/internal/repository"
)

type BannerService struct {
	repo repository.BannerStorage
}

func NewBannerService(repo repository.BannerStorage) *BannerService {
	return &BannerService{repo: repo}
}

func (b *BannerService) GetUserBannerAction(ctx context.Context, params GetUserBannerParameters) (models.UserBanner, error) {
	return models.UserBanner{}, nil
}

func (b *BannerService) GetBannerWithFilterAction(ctx context.Context, params GetBannerWithFilterParameters) (models.UserBanner, error) {
	return models.UserBanner{}, nil
}
