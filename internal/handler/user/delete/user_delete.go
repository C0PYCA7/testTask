package delete

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	"testTask/internal/lib/response"
)

type Response struct {
	response.Response
}

type DeleteUser interface {
	DeleteUser(id int64) error
}

func New(log *slog.Logger, deleteUser DeleteUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler/user/delete/DeleteHandler"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userIdStr := chi.URLParam(r, "id")

		if userIdStr != "" {
			userId, err := strconv.ParseInt(userIdStr, 10, 64)
			if err != nil {
				log.Error("invalid user ID ", err)

				render.JSON(w, r, "invalid user id")

				return
			}
			log.Info("got user id ", slog.Int64("id", userId))

			err = deleteUser.DeleteUser(userId)
			if err != nil {
				log.Error("failed to delete user ", err)

				render.JSON(w, r, "failed to delete user")

				return
			}

			log.Info("user delete")

			render.JSON(w, r, response.OK())

		} else {
			log.Error("empty user ID")

			render.JSON(w, r, "empty user ID")

			return
		}
	}
}
