package create

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"testTask/internal/lib/response"
)

type Request struct {
	Name       string `json:"name" validate:"required"`
	Surname    string `json:"surname" validate:"required"`
	Patronymic string `json:"patronymic,omitempty"`
}

type Response struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
	response.Response
}

type AgeResponse struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

type GenderResponse struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}

type NationalityResponse struct {
	Count   int             `json:"count"`
	Name    string          `json:"name"`
	Country []CountryDetail `json:"country"`
}

type CountryDetail struct {
	//NationalityMap map[string]float64 `json:"nationality_map"`
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type CreateUser interface {
	CreateUser(name, surname, patronymic, nationality string, age int, gender bool) (int64, error)
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers/user/user_create/New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body: ", err)

			render.JSON(w, r, response.Error("failed to decode request body"))

			return
		}

		log.Info("request body decoded")

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request: ", err)

			render.JSON(w, r, "invalid request")

			return
		}

		userAge, err := enrichUserAge(&req)
		if err != nil {
			log.Error("failed to get user age")

			render.JSON(w, r, "failed to get user data")

			return
		}

		userGender, err := enrichUserGender(&req)
		if err != nil {
			log.Error("failed to det user gender")

			render.JSON(w, r, "failed to get user gender")

			return
		}

		userNationality, err := enrichUserNationality(&req)
		if err != nil {
			log.Error("failed to det user gender")

			render.JSON(w, r, "failed to get user gender")

			return
		}

		render.JSON(w, r, Response{
			Name:        req.Name,
			Surname:     req.Surname,
			Patronymic:  req.Patronymic,
			Age:         userAge.Age,
			Gender:      userGender.Gender,
			Nationality: GetMaxProbabilityNationality(userNationality),
			Response:    response.OK(),
		})
	}
}

func enrichUserAge(req *Request) (*AgeResponse, error) {

	const op = "user/create/user_create/enrichUserAge"

	ageURL := fmt.Sprintf("https://api.agify.io/?name=%s", req.Name)
	agifyResp, err := http.Get(ageURL)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get age data: %w", op, err)
	}
	defer agifyResp.Body.Close()

	var ageData AgeResponse

	if err := render.DecodeJSON(agifyResp.Body, &ageData); err != nil {
		return nil, fmt.Errorf("%s: failed to decode response body: %w", op, err)
	}
	return &ageData, nil
}

func enrichUserGender(req *Request) (*GenderResponse, error) {
	const op = "user/create/user_create/enrichUserGender"

	genderUrl := fmt.Sprintf("https://api.genderize.io/?name=%s", req.Name)
	agifyResp, err := http.Get(genderUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get gender data: %w", op, err)
	}
	defer agifyResp.Body.Close()

	var genderData GenderResponse

	if err := render.DecodeJSON(agifyResp.Body, &genderData); err != nil {
		return nil, fmt.Errorf("%s: failed to decode response body: %w", op, err)
	}
	return &genderData, nil
}

func enrichUserNationality(req *Request) (*NationalityResponse, error) {
	const op = "user/create/user_create/enrichUserNationality"

	nationalityUrl := fmt.Sprintf("https://api.nationalize.io/?name=%s", req.Name)

	agifyResp, err := http.Get(nationalityUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get nationaluty data: %w", op, err)
	}
	defer agifyResp.Body.Close()

	var nationalityData NationalityResponse

	if err := render.DecodeJSON(agifyResp.Body, &nationalityData); err != nil {
		return nil, fmt.Errorf("%s: failed to decode response body: %w", op, err)
	}
	return &nationalityData, nil
}

func GetMaxProbabilityNationality(nationalResponse *NationalityResponse) string {

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
