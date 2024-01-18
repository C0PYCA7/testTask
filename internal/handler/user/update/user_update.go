package update

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
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
	UpdateUser(id int64, request Request) (int64, error)
}

func New(log *slog.Logger, userUpdate UserUpdate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler/user/update/user_update/New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body ", err)

			render.JSON(w, r, "failed to decode request body")

			return
		}

		userIdStr := chi.URLParam(r, "id")
		if userIdStr != "" {
			userId, err := strconv.ParseInt(userIdStr, 10, 64)
			if err != nil {
				log.Error("invalid user ID ", err)

				render.JSON(w, r, "invalid user id")

				return
			}

			id, err := userUpdate.UpdateUser(userId, req)
			if err != nil {
				log.Error("failed to update user: ", err)

				render.JSON(w, r, "failed to update user")

				return
			}

			log.Info("update user: ", id)

			render.JSON(w, r, response.OK())
		} else {
			log.Error("empty user ID")

			render.JSON(w, r, "empty user ID")

			return
		}
	}
}
