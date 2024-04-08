package models

type CreateBanner struct {
	TagID     []string
	FeatureID int
	NewBanner struct{}
	IsActive  bool
}

type GetAllBanners struct {
	TagID     string
	FeatureID string
	Limit     string
	Offset    string
}
