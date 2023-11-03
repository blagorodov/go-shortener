package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type claims struct {
	jwt.RegisteredClaims
	UserID int
}

const (
	TOKEN_EXP  = time.Hour * 3
	SECRET_KEY = "secretsecretsecret"
)

func EncodeToken(id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: id,
	})

	tokenString, err := token.SignedString([]byte(SECRET_KEY))
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
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, fmt.Errorf("token is not valid")
	}

	return cls.UserID, nil
}
