package banners

import (
	response "banner/internal/lib/api/responses"
	logerr "banner/internal/lib/logger/logerr"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type RequestUpdateBanner struct {
	TagIDs    []int                  `json:"tag_ids" validate:"required"`
	FeatureID int                    `json:"feature_id" validate:"required"`
	Content   map[string]interface{} `json:"content" validate:"required"`
	IsActive  bool                   `json:"is_active" validate:"required"`
}

func UpdateBanner(bannerRepo Banners, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bannerID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			logger.Error("Invalid banner ID", logerr.Err(err))
			render.JSON(w, r, response.Error("Invalid banner ID"))
			return
		}

		var req RequestUpdateBanner
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			render.Status(r, http.StatusBadRequest)
			logger.Error("Failed to parse request body", logerr.Err(err))
			render.JSON(w, r, response.Error("Failed to parse request body"))
			return
		}

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			logger.Error("Invalid request", logerr.Err(err))
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}

		banner, err := bannerRepo.FindBannerId(r.Context(), bannerID)
		if err != nil {
			render.Status(r, http.StatusNotFound)
			logger.Error("Failed to find banner")
			render.JSON(w, r, response.Error("Banner not found"))
			return
		}

		banner.TagIDs = req.TagIDs
		banner.FeatureID = req.FeatureID
		banner.Content = req.Content
		banner.IsActive = req.IsActive
		banner.UpdatedAt = time.Now()

		if err := bannerRepo.UpdateBanner(r.Context(), &banner); err != nil {
			render.Status(r, http.StatusInternalServerError)
			logger.Error("Failed to update banner", logerr.Err(err))
			render.JSON(w, r, response.Error("Failed to update banner"))
			return
		}

		ResponseOK(w, r, banner)
	}
}
