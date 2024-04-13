package banners

import (
	response "banner/internal/lib/api/responses"
	logerr "banner/internal/lib/logger/logerr"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func DeleteBanner(log *slog.Logger, bannerRepo Banners) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.banners.deleteBanner.New"
		log := log.With(
			slog.String("options", loggerOptions),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("Invalid banner ID"))
			return
		}

		err = bannerRepo.DeleteBannerID(r.Context(), id)
		if err != nil {
			log.Error("Failed to delete banner", logerr.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("Failed to delete banner"))
			return
		}

		log.Info("Banner deleted")
		render.Status(r, http.StatusNoContent)
	}
}
