package tags

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

type RequestTag struct {
	Name string `json:"name" validate:"required"`
}
type ResponseTag struct {
	response.Response
	ID   int    `json:"tag_id"`
	Name string `json:"name"`
}

type Tag interface {
	CreateTag(ctx context.Context, tag *models.Tag) error
}

func NewTag(log *slog.Logger, tagRepo Tag) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.features.createTag.New"
		log = log.With(
			slog.String("options", loggerOptions),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req RequestTag
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

		tag := models.Tag{Name: req.Name}
		err = tagRepo.CreateTag(r.Context(), &tag)
		if err != nil {
			log.Error("Failed to create tag", logerr.Err(err))
			render.JSON(w, r, response.Error("Failed to create tag"))
			return
		}

		log.Info("Tag added")
		ResponseOK(w, r, req.Name, tag.ID)
	}
}

func ResponseOK(w http.ResponseWriter, r *http.Request, name string, tag_id int) {
	render.JSON(w, r, ResponseTag{Response: response.OK(), Name: name, ID: tag_id})
}
