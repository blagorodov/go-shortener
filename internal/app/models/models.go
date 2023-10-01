package models

type (
	ShortenRequest struct {
		Url string `json:"url"`
	}

	ShortenResponse struct {
		Result string `json:"result"`
	}
)
