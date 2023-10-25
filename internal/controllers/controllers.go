package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/blagorodov/go-shortener/internal/config"
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

	url, err := parseOne(r)
	if err != nil {
		return "", err
	}

	if len(url) == 0 {
		return "", errors.New("пустой url")
	}

	key, err := s.NewKey(ctx)
	if err != nil {
		return "", err
	}

	err = s.Put(ctx, key, url)
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

		err = s.Put(ctx, key, item.OriginalURL)
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
		return readBody(r), nil
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

func readBody(r *http.Request) string {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(r.Body)
	body, _ := io.ReadAll(r.Body)
	return string(body)
}
