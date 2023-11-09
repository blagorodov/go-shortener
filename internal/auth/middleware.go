package auth

import (
	"github.com/blagorodov/go-shortener/internal/cookies"
	"net/http"
)

func TokenMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// получим токен из куки, если нет или невалидный, то сгенерим токен
		if !cookies.Check(r) {
			// установим токен в куку
			http.SetCookie(w, cookies.New())
		}

		// передаём управление хендлеру
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
