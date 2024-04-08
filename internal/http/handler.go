package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"banner/internal/handler"
	"banner/internal/repository"
)

type wrappedHandler func(w http.ResponseWriter, r *http.Request) error

type Handler struct {
	service *handler.BannerService
}

func NewHandler(service *handler.BannerService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) getUserBanner(w http.ResponseWriter, r *http.Request) error {
	tagID := r.URL.Query().Get("tag_id")
	featureID := r.URL.Query().Get("feature_id")
	useLastRevision := r.URL.Query().Get("use_last_revision")

	result, err := h.service.GetUserBannerAction(r.Context(), handler.GetUserBannerParameters{
		TagID:           tagID,
		FeatureID:       featureID,
		UseLastRevision: useLastRevision == "true",
	})
	if err != nil {
		return err
	}

	return sendJSONResponse(w, result, http.StatusOK)
}

func (h *Handler) getBannerWithFilter(w http.ResponseWriter, r *http.Request) error {
	tagID := r.URL.Query().Get("tag_id")
	featureID := r.URL.Query().Get("feature_id")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	result, err := h.service.GetBannerWithFilterAction(r.Context(), handler.GetBannerWithFilterParameters{
		TagID:     tagID,
		FeatureID: featureID,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return err
	}

	return sendJSONResponse(w, result, http.StatusOK)
}

func (h *Handler) createBanner(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *Handler) updateBanner(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *Handler) deleteBanner(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func errorsHandler(handler wrappedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)

		switch {
		case err == nil:
			w.WriteHeader(http.StatusOK)
		case errors.Is(err, repository.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
		default:
			msg := fmt.Sprintf(`{"error":"%s"}`, err.Error())
			_ = sendJSONResponse(w, msg, http.StatusInternalServerError)
		}
	}
}

func sendJSONResponse(w http.ResponseWriter, result interface{}, status int) error {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		return err
	}

	return nil
}
