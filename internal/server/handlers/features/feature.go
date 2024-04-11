package features

import (
	response "banner/internal/lib/api/responses"
	logerr "banner/internal/lib/logger/logerr"
	"banner/internal/models"
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type RequestFeature struct {
	Name string `json:"name" validate:"required"`
}

type ResponseFeature struct {
	response.Response
	ID   int    `json:"feature_id"`
	Name string `json:"name"`
}

type Features interface {
	CreateFeature(ctx context.Context, feature *models.Feature) error
}

func NewFeature(log *slog.Logger, featureRepo Features) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.features.createFeature.New"
		log = log.With(
			slog.String("options", loggerOptions),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req RequestFeature
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to decode request body", logerr.Err(err))
			render.JSON(w, r, response.Error("Failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("Invalid request", logerr.Err(err))
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}

		feature := models.Feature{Name: req.Name}
		err = featureRepo.CreateFeature(r.Context(), &feature)
		if err != nil {
			log.Error("Failed to create feature", logerr.Err(err))
			render.JSON(w, r, response.Error("Failed to create feature"))
			return
		}

		log.Info("Feature added")
		ResponseOK(w, r, req.Name, feature.ID)
	}
}

func ResponseOK(w http.ResponseWriter, r *http.Request, name string, feature_id int) {
	render.JSON(w, r, ResponseFeature{Response: response.OK(),
		Name: name, ID: feature_id})
}
