package server

import (
	"github.com/blagorodov/go-shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, path, body)
	require.NoError(t, err)

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	ts := httptest.NewServer(Router())
	defer ts.Close()

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
				s := strings.TrimPrefix(savedLink, "http://example.com")
				route = s
				requestBody = ""
			} else {
				route = ts.URL
			}

			resp, respBody := testRequest(t, ts, tc.method, route, strings.NewReader(requestBody))

			assert.Equal(t, tc.expectedCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")

			t.Log(resp.Header.Values("Location"))

			if tc.saveResult {
				savedLink = respBody
			}
			if tc.expectedHeaderKey != "" {
				assert.Equal(t, tc.expectedHeaderValue, resp.Header.Get(tc.expectedHeaderKey), "Заголовок не совпадает с ожиданием")
			}
		})
	}

}
