package update

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

type Request struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic,omitempty"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
	Age         int    `json:"age"`
}

type Response struct {
	ID int64 `json:"id"`
	response.Response
}

type UserUpdate interface {
	UpdateUser(id int64, request Request) error
}

func New(log *slog.Logger, userUpdate UserUpdate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler/user/update/user_update/New"

		log = log.With(
			slog.String("op", op),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", err)

			render.JSON(w, r, "failed to decode request body")

			return
		}

		log.Debug("request data from user", slog.Any("req", req))

		userIdStr := chi.URLParam(r, "id")
		if userIdStr != "" {
			userId, err := strconv.ParseInt(userIdStr, 10, 64)
			if err != nil {
				log.Error("invalid user ID", err)

				render.JSON(w, r, "invalid user id")

				return
			}

			log.Debug("userID", slog.Any("Id", userId))

			err = userUpdate.UpdateUser(userId, req)
			if err != nil {
				if errors.Is(err, database.ErrUserNotFound) {
					log.Error("user not found")

					w.WriteHeader(http.StatusNotFound)

					render.JSON(w, r, "user not found")

					return
				}
				log.Error("failed to update user", err)

				w.WriteHeader(http.StatusInternalServerError)

				render.JSON(w, r, "failed to update user")

				return
			}

			log.Info("update user", slog.Int64("id", userId))

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
