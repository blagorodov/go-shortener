package handlers

import (
	"github.com/blagorodov/go-shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRoot(t *testing.T) {
	storage.Init()

	testCases := []struct {
		method              string
		expectedCode        int
		requestBody         string
		saveResult          bool
		expectedHeaderKey   string
		expectedHeaderValue string
	}{
		{
			method:              http.MethodPost,
			expectedCode:        http.StatusCreated,
			requestBody:         "https://practicum.yandex.ru/",
			saveResult:          true,
			expectedHeaderKey:   "",
			expectedHeaderValue: "",
		},
		{
			method:              http.MethodGet,
			expectedCode:        http.StatusTemporaryRedirect,
			requestBody:         "",
			saveResult:          false,
			expectedHeaderKey:   "Location",
			expectedHeaderValue: "https://practicum.yandex.ru/",
		},
	}

	var savedLink string

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			requestBody := tc.requestBody
			route := "/"
			if !tc.saveResult {
				s, _ := strings.CutPrefix(savedLink, "http://example.com")
				route = s
				requestBody = ""
			}

			r := httptest.NewRequest(tc.method, route, strings.NewReader(requestBody))
			w := httptest.NewRecorder()

			Root(w, r)

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")

			t.Log(w.Header().Values("Location"))

			if tc.saveResult {
				savedLink = w.Body.String()
			}
			if tc.expectedHeaderKey != "" {
				assert.Equal(t, tc.expectedHeaderValue, w.Header().Get(tc.expectedHeaderKey), "Заголовок не совпадает с ожиданием")
			}
		})
	}
}
