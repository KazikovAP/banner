package users

import (
	response "banner/internal/lib/api/responses"
	password "banner/internal/lib/auth/password"
	logerr "banner/internal/lib/logger/logerr"
	"banner/internal/models"
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type RequestUser struct {
	Username string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ResponseUser struct {
	response.Response
	ID   int    `json:"user_id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type User interface {
	CreateUser(ctx context.Context, user *models.User) error
	FindUserUsername(ctx context.Context, username string) (models.User, error)
}

func NewUser(log *slog.Logger, u User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.features.createUser.New"
		log = log.With(
			slog.String("options", loggerOptions),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req RequestUser
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

		hashPass, err := password.HashPassword(req.Password)
		user := models.User{Username: req.Username, Password: hashPass, Role: "user"}
		err = u.CreateUser(r.Context(), &user)
		if err != nil {
			log.Error("Failed to create user", logerr.Err(err))
			render.JSON(w, r, response.Error("Failed to create user"))
			return
		}

		log.Info("User added")
		ResponseOK(w, r, req.Username, user.ID, user.Role)
	}
}

func ResponseOK(w http.ResponseWriter, r *http.Request, name string, userID int, role string) {
	render.JSON(w, r, ResponseUser{Response: response.OK(),
		Name: name, ID: userID, Role: role})
}
