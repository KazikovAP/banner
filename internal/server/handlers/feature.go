package handlers

import (
	response "banner/internal/lib/api"
	"banner/internal/lib/logger"
	"banner/internal/models"
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Name string `json:"name" validate:"required"`
}

type Response struct {
	response.Response
	ID   int    `json:"feature_id"`
	Name string `json:"name"`
}

type Features interface {
	CreateFeature(ctx context.Context, feature *models.Feature) error
}

func NewFeature(log *slog.Logger, feat Features) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.features.createFeature.NewFeature"

		log = log.With(
			slog.String("options", loggerOptions),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req Request
		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("Failed to decode request body", logger.Err(err))
			render.JSON(w, r, response.Error("Failed to decode request"))
			return
		}

		log.Info("Request body decoded", slog.Any("Request", req))
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("Invalid request", logger.Err(err))
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}

		feature := models.Feature{Name: req.Name}
		err = feat.CreateFeature(r.Context(), &feature)
		if err != nil {
			log.Error("Failed to create feature", logger.Err(err))
			render.JSON(w, r, response.Error("Failed to create feature"))
			return
		}

		log.Info("Feature added")
		ResponseOK(w, r, req.Name)
	}
}

func ResponseOK(w http.ResponseWriter, r *http.Request, name string) {
	render.JSON(w, r, Response{Response: response.OK(), Name: name})
}
