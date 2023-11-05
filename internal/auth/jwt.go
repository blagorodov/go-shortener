package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

type claims struct {
	jwt.RegisteredClaims
	UserID int
}

const (
	TokenExp  = time.Hour * 3
	SecretKey = "secretsecretsecret"
)

func EncodeToken(id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: id,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func DecodeToken(tokenString string) (int, error) {
	cls := &claims{}

	token, err := jwt.ParseWithClaims(tokenString, cls, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(SecretKey), nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, fmt.Errorf("token is not valid")
	}

	return cls.UserID, nil
}

func GetUserID(r *http.Request) (int, error) {
	cookie, err := r.Cookie("x-token")
	if err != nil {
		return 0, err
	}
	return DecodeToken(cookie.Value)
}
