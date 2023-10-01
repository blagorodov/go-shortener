package server

import (
	"github.com/blagorodov/go-shortener/internal/app/config"
	"github.com/blagorodov/go-shortener/internal/app/logger"
	"github.com/blagorodov/go-shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, contentType, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", contentType)
	require.NoError(t, err)

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	logger.Init()
	ts := httptest.NewServer(router(storage.NewMemoryStorage()))
	defer ts.Close()

	testCases := []struct {
		route               string
		method              string
		contentType         string
		expectedCode        int
		requestBody         string
		saveResult          bool
		expectedHeaderKey   string
		expectedHeaderValue string
	}{
		{
			route:               "/",
			method:              http.MethodPost,
			contentType:         "text/plain",
			expectedCode:        http.StatusCreated,
			requestBody:         "https://practicum.yandex.ru/",
			saveResult:          true,
			expectedHeaderKey:   "",
			expectedHeaderValue: "",
		},
		{
			route:               "/",
			method:              http.MethodGet,
			contentType:         "text/plain",
			expectedCode:        http.StatusTemporaryRedirect,
			requestBody:         "",
			saveResult:          false,
			expectedHeaderKey:   "Location",
			expectedHeaderValue: "https://practicum.yandex.ru/",
		},
		{
			route:               "/api/shorten",
			method:              http.MethodPost,
			contentType:         "application/json",
			expectedCode:        http.StatusCreated,
			requestBody:         `{"url":"https://practicum.yandex.ru/"}`,
			saveResult:          true,
			expectedHeaderKey:   "Content-Type",
			expectedHeaderValue: "application/json",
		},
	}

	var savedLink string

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			requestBody := tc.requestBody
			route := tc.route
			if !tc.saveResult {
				s := strings.TrimPrefix(savedLink, config.Options.BaseURL)
				route = ts.URL + s
				requestBody = ""
			} else {
				route = ts.URL
			}
			resp, respBody := testRequest(t, ts, tc.method, tc.contentType, route, strings.NewReader(requestBody))

			assert.Equal(t, tc.expectedCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")

			if tc.saveResult {
				savedLink = respBody
			}
			if tc.expectedHeaderKey != "" {
				assert.Equal(t, tc.expectedHeaderValue, resp.Header.Get(tc.expectedHeaderKey), "Заголовок не совпадает с ожиданием")
			}

			err := resp.Body.Close()
			if err != nil {
				return
			}
		})
	}

}
