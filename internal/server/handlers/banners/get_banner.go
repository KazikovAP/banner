package banners

import (
	response "banner/internal/lib/api/responses"
	logerr "banner/internal/lib/logger/logerr"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
)

type RequestGetBanners struct {
	FeatureID *int `json:"feature_id"`
	TagID     *int `json:"tag_id"`
	Limit     *int `json:"limit"`
	Offset    *int `json:"offset"`
}

func GetBanners(bannerRepo Banners, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := ParseGetBannersRequest(r)

		banners, err := bannerRepo.FindBannersParameters(r.Context(), req)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			logger.Error("Failed to get banner", logerr.Err(err))
			render.JSON(w, r, response.Error("Failed to get banner"))
			return
		}

		render.JSON(w, r, banners)
	}
}

func ParseGetBannersRequest(r *http.Request) RequestGetBanners {
	req := RequestGetBanners{}

	if featureIDStr := r.URL.Query().Get("feature_id"); featureIDStr != "" {
		featureID, _ := strconv.Atoi(featureIDStr)
		req.FeatureID = &featureID
	}

	if tagIDStr := r.URL.Query().Get("tag_id"); tagIDStr != "" {
		tagID, _ := strconv.Atoi(tagIDStr)
		req.TagID = &tagID
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, _ := strconv.Atoi(limitStr)
		req.Limit = &limit
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, _ := strconv.Atoi(offsetStr)
		req.Offset = &offset
	}

	return req
}
