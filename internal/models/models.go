package models

type Filter struct {
	Name        string
	Surname     string
	Patronymic  string
	Age         int
	Gender      string
	Nationality string
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
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}
