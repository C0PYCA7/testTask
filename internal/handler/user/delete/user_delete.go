package delete

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	"testTask/internal/database"
	"testTask/internal/lib/response"
)

type Response struct {
	ID int64 `json:"userID"`
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
		)

		userIdStr := chi.URLParam(r, "id")

		if userIdStr != "" {
			userId, err := strconv.ParseInt(userIdStr, 10, 64)
			if err != nil {
				log.Error("invalid user ID", err)

				render.JSON(w, r, "invalid user id")

				return
			}
			log.Debug("userID", slog.Any("id", userId))

			err = deleteUser.DeleteUser(userId)
			if err != nil {
				if errors.Is(err, database.ErrUserNotFound) {

					w.WriteHeader(http.StatusNotFound)

					log.Error("user not found")

					render.JSON(w, r, "user not found")

					return
				}
				log.Error("failed to delete user", err)

				w.WriteHeader(http.StatusInternalServerError)

				render.JSON(w, r, "failed to delete user")

				return
			}

			log.Info("user delete")

			render.JSON(w, r, Response{
				ID:       userId,
				Response: response.OK(),
			})

		} else {
			log.Error("empty user ID")

			render.JSON(w, r, "empty user ID")

			return
		}
	}
}
