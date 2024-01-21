package create

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"testTask/internal/database/postgres"
	"testTask/internal/lib/response"
	"testTask/internal/models"
)

type Request struct {
	Name       string `json:"name" validate:"required"`
	Surname    string `json:"surname" validate:"required"`
	Patronymic string `json:"patronymic,omitempty"`
}

type Response struct {
	Id int64 `json:"id"`
	response.Response
}

type CreateUser interface {
	CreateUser(user *postgres.User) (int64, error)
}

func New(log *slog.Logger, createUser CreateUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers/user/user_create/New"

		log = log.With(
			slog.String("op", op),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body: ", err)

			render.JSON(w, r, "failed to decode request body")

			return
		}

		log.Info("request body decoded")

		log.Debug("request data from user", slog.Any("req", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request: ", err)

			render.JSON(w, r, "invalid request")

			return
		}

		ageCh := make(chan *models.AgeResponse)
		genderCh := make(chan *models.GenderResponse)
		nationalityCh := make(chan *models.NationalityResponse)

		go func() {
			ageData, err := enrichUserAge(&req)
			if err != nil {
				log.Error("failed to get user age")
				ageCh <- nil
			} else {
				ageCh <- ageData
			}
		}()

		go func() {
			genderData, err := enrichUserGender(&req)
			if err != nil {
				log.Error("failed to get user gender")
				genderCh <- nil
			} else {
				genderCh <- genderData
			}
		}()

		go func() {
			nationalityData, err := enrichUserNationality(&req)
			if err != nil {
				log.Error("failed to get user nationality")
				nationalityCh <- nil
			} else {
				nationalityCh <- nationalityData
			}
		}()

		userAge := <-ageCh
		userGender := <-genderCh
		userNationality := <-nationalityCh

		nationality := GetMaxProbabilityNationality(userNationality)
		if nationality == "" || userGender == nil || userAge == nil {
			log.Info("empty user data")

			render.JSON(w, r, "empty data")

			return
		}

		log.Debug("data from other API", slog.Any("age", userAge), slog.Any("gender", userGender), slog.Any("nationality", nationality))

		user := postgres.User{
			Name:        req.Name,
			Surname:     req.Surname,
			Patronymic:  req.Patronymic,
			Age:         userAge.Age,
			Gender:      userGender.Gender,
			Nationality: nationality,
		}

		id, err := createUser.CreateUser(&user)
		if err != nil {
			log.Error("failed to add user", err)

			render.JSON(w, r, "failed to add user")

			return
		}

		log.Info("user added", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Id:       id,
			Response: response.OK(),
		})
	}
}

func enrichUserAge(req *Request) (*models.AgeResponse, error) {

	const op = "user/create/user_create/enrichUserAge"

	ageURL := fmt.Sprintf("https://api.agify.io/?name=%s", req.Name)
	agifyResp, err := http.Get(ageURL)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get age data: %w", op, err)
	}
	defer agifyResp.Body.Close()

	var ageData models.AgeResponse

	if err := render.DecodeJSON(agifyResp.Body, &ageData); err != nil {
		return nil, fmt.Errorf("%s: failed to decode response body: %w", op, err)
	}
	return &ageData, nil
}

func enrichUserGender(req *Request) (*models.GenderResponse, error) {
	const op = "user/create/user_create/enrichUserGender"

	genderUrl := fmt.Sprintf("https://api.genderize.io/?name=%s", req.Name)
	agifyResp, err := http.Get(genderUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get gender data: %w", op, err)
	}
	defer agifyResp.Body.Close()

	var genderData models.GenderResponse

	if err := render.DecodeJSON(agifyResp.Body, &genderData); err != nil {
		return nil, fmt.Errorf("%s: failed to decode response body: %w", op, err)
	}
	return &genderData, nil
}

func enrichUserNationality(req *Request) (*models.NationalityResponse, error) {
	const op = "user/create/user_create/enrichUserNationality"

	nationalityUrl := fmt.Sprintf("https://api.nationalize.io/?name=%s", req.Name)

	agifyResp, err := http.Get(nationalityUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get nationaluty data: %w", op, err)
	}
	defer agifyResp.Body.Close()

	var nationalityData models.NationalityResponse

	if err := render.DecodeJSON(agifyResp.Body, &nationalityData); err != nil {
		return nil, fmt.Errorf("%s: failed to decode response body: %w", op, err)
	}
	return &nationalityData, nil
}

func GetMaxProbabilityNationality(nationalResponse *models.NationalityResponse) string {

	if len(nationalResponse.Country) == 0 {
		return ""
	}

	maxProbability := nationalResponse.Country[0].Probability
	maxCountry := nationalResponse.Country[0].CountryID

	for _, country := range nationalResponse.Country {
		if country.Probability > maxProbability {
			maxProbability = country.Probability
			maxCountry = country.CountryID
		}
	}
	return maxCountry
}
