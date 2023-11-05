package auth

import (
	"errors"
	"math/rand"
	"net/http"
	"time"
)

func TokenMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var token string
		createNew := false

		// получим токен из куки, если нет или невалидный, то сгенерим токен
		cookie, err := r.Cookie("x-token")
		if errors.Is(err, http.ErrNoCookie) {
			createNew = true
		} else {
			token = cookie.Value
			if _, err := DecodeToken(token); err != nil {
				createNew = true
			}
		}

		// сгенерим новый токен если требуется
		if createNew {
			token, err = EncodeToken(rand.Intn(1000000))
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		// установим токен в куку
		http.SetCookie(w, &http.Cookie{
			Name:   "x-token",
			Value:  token,
			MaxAge: (int)(time.Hour * 24 * 30),
		})

		// передаём управление хендлеру
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
