package login

import (
	response "banner/internal/lib/api/responses"
	jwt "banner/internal/lib/auth/jwt"
	password "banner/internal/lib/auth/password"
	logerr "banner/internal/lib/logger/logerr"
	user "banner/internal/server/handlers/users/user"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type ResponseAuthUser struct {
	response.Response
	ID    int    `json:"user_id"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	Token string `json:"token"`
}

func Login(log *slog.Logger, userRepo user.User, jwt *jwt.JWTSecret) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.features.createUser.New"
		log = log.With(
			slog.String("options", loggerOptions),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req user.RequestUser
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

		user, err := userRepo.FindUserUsername(r.Context(), req.Username)
		if err != nil {
			log.Error("User not found with login")
			render.JSON(w, r, response.Error("Invalid username"))
			return
		}

		errAuth := password.ComparePasswordHash(req.Password, user.Password)
		if errAuth != nil {
			log.Error("Invalid password")
			render.JSON(w, r, response.Error("Invalid password"))
			return
		}

		token, err := jwt.GenerateToken(user.Username, user.Role, time.Second*600)
		log.Info("User authenticated")
		ResponseAuthOK(w, r, req.Username, user.ID, user.Role, token)
	}
}

func ResponseAuthOK(w http.ResponseWriter, r *http.Request, name string, userID int, role string, token string) {
	render.JSON(w, r, ResponseAuthUser{Response: response.OK(),
		Name: name, ID: userID, Role: role, Token: token})
}
