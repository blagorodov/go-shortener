package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blagorodov/go-shortener/internal/config"
	"github.com/blagorodov/go-shortener/internal/cookies"
	"github.com/blagorodov/go-shortener/internal/errs"
	"github.com/blagorodov/go-shortener/internal/models"
	"github.com/blagorodov/go-shortener/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"io"
	"net/http"
	"strings"
)

// Get Контроллер GET /
func Get(ctx context.Context, r *http.Request, s service.Service) (string, error) {
	return s.Get(ctx, chi.URLParam(r, "id"))
}

// ShortenOne Контроллер POST /api/shorten
func ShortenOne(ctx context.Context, r *http.Request, s service.Service) (string, error) {
	var url string
	var resultErr error
	fmt.Println("shortone")

	url, err := parseOne(r)
	if err != nil {
		return "", err
	}

	if len(url) == 0 {
		return "", errors.New(errs.ErrEmptyURL)
	}

	key, err := s.NewKey(ctx)
	if err != nil {
		return "", err
	}

	userID, _ := cookies.GetID(r)
	fmt.Println(userID)

	err = s.Put(ctx, key, url, userID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		var errKey error
		key, errKey = s.GetKey(ctx, url)
		if errKey != nil {
			return "", errKey
		}
		resultErr = errors.New(errs.ErrUniqueLinkCode)
	} else if err != nil {
		return "", err
	}

	parts := []string{config.Options.BaseURL, key}
	url = strings.Join(parts, `/`)

	return url, resultErr
}

// ShortenBatch Контроллер POST /api/shorten/batch
func ShortenBatch(ctx context.Context, r *http.Request, s service.Service) (models.BatchResponseList, error) {
	var urls models.BatchRequestList
	var result models.BatchResponseList
	var resultErr error

	urls, err := parseBatch(r)
	if err != nil {
		return nil, err
	}

	if len(urls) == 0 {
		return nil, nil
	}

	for _, item := range urls {
		key, err := s.NewKey(ctx)
		if err != nil {
			return nil, err
		}

		userID, _ := cookies.GetID(r)

		err = s.Put(ctx, key, item.OriginalURL, userID)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			var errKey error
			key, errKey = s.GetKey(ctx, item.OriginalURL)
			if errKey != nil {
				return nil, errKey
			}
			resultErr = errors.New(errs.ErrUniqueLinkCode)
		} else if err != nil {
			return nil, err
		}

		parts := []string{config.Options.BaseURL, key}
		url := strings.Join(parts, `/`)

		result = append(result, models.BatchResponse{
			CorrelationID: item.CorrelationID,
			ShortURL:      url,
		})
	}

	return result, resultErr
}

func GetURLs(ctx context.Context, r *http.Request, s service.Service) (models.AllResponseList, error) {
	userID, err := cookies.GetID(r)
	urls, errURLs := s.GetURLs(ctx, userID)
	if errURLs != nil {
		err = errURLs
	}
	return urls, err
}

func parseOne(r *http.Request) (string, error) {
	if r.Header.Get("Content-Type") == "application/json" {
		var request models.ShortenRequest
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			return "", err
		}
		if err = json.Unmarshal(buf.Bytes(), &request); err != nil {
			return "", err
		}

		return request.URL, nil
	} else {
		return readBody(r)
	}
}

func parseBatch(r *http.Request) (models.BatchRequestList, error) {
	if r.Header.Get("Content-Type") == "application/json" {
		var list models.BatchRequestList
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			return nil, err
		}
		if err = json.Unmarshal(buf.Bytes(), &list); err != nil {
			return nil, err
		}

		return list, nil
	}
	return nil, nil
}

func parseLogin(r *http.Request) (*models.LoginRequest, error) {
	if r.Header.Get("Content-Type") == "application/json" {
		var buf bytes.Buffer
		var result models.LoginRequest

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			return nil, err
		}
		if err = json.Unmarshal(buf.Bytes(), &result); err != nil {
			return nil, err
		}

		return &result, nil
	}
	return nil, nil
}

func readBody(r *http.Request) (string, error) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(r.Body)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
