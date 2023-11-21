package models

type (
	ShortenRequest struct {
		URL string `json:"url"`
	}

	ShortenResponse struct {
		Result string `json:"result"`
	}

	BatchRequest struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	BatchResponse struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}

	AllResponse struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	BatchRequestList []BatchRequest

	BatchResponseList []BatchResponse

	AllResponseList []AllResponse

	LoginRequest struct {
		UserID int `json:"user_id"`
	}

	LoginResponse struct {
		Token string `json:"token"`
	}
)
