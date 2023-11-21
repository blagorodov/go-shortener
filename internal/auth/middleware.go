package auth

import (
	"context"
	"github.com/blagorodov/go-shortener/internal/cookies"
	"net/http"
)

type key int

const ContextKey key = 1

func TokenMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// получим токен из куки, если нет или невалидный, то сгенерим токен
		userID, _ := cookies.GetID(w, r)
		r = r.WithContext(context.WithValue(r.Context(), ContextKey, userID))

		// передаём управление хендлеру
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
