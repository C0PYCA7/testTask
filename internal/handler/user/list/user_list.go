package list

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	"testTask/internal/database/postgres"
	"testTask/internal/lib/response"
	"testTask/internal/models"
)

//type Request struct {
//	Name        string `json:"name,omitempty"`
//	Surname     string `json:"surname,omitempty"`
//	Patronymic  string `json:"patronymic,omitempty"`
//	Age         int    `json:"age,omitempty"`
//	Gender      string `json:"gender,omitempty"`
//	Nationality string `json:"nationality,omitempty"`
//}

type Response struct {
	Users []postgres.User `json:"users"`
	response.Response
}

type GetUsers interface {
	GetUsers(filter models.Filter, pageSize, page int) ([]postgres.User, error)
}

func New(log *slog.Logger, getUsers GetUsers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler/user/list/user_list/New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var age int
		ageStr := r.URL.Query().Get("age")
		if ageStr != "" {
			var err error
			age, err = strconv.Atoi(ageStr)
			if err != nil {
				log.Error("failed to parse age")

				render.JSON(w, r, "failed to parse age")

				return
			}
		}

		filter := models.Filter{
			Name:        r.URL.Query().Get("name"),
			Surname:     r.URL.Query().Get("surname"),
			Patronymic:  r.URL.Query().Get("patronymic"),
			Age:         age,
			Gender:      r.URL.Query().Get("gender"),
			Nationality: r.URL.Query().Get("nationality"),
		}

		pageSizeStr := r.URL.Query().Get("size")
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			log.Error("failed to parse size")

			render.JSON(w, r, "failed to parse size")

			return
		}
		pageStr := r.URL.Query().Get("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			log.Error("failed to parse page")

			render.JSON(w, r, "failed to parse page")

			return
		}

		users, err := getUsers.GetUsers(filter, pageSize, page)
		if err != nil {
			log.Error("failed to get users")

			render.JSON(w, r, "failed to get users")

			return
		}

		log.Info("got users")

		render.JSON(w, r, Response{users, response.OK()})
	}
}
