package cookies

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

var hashKey []byte

func New() *http.Cookie {
	id := uuid.NewString()

	h := hmac.New(sha256.New, hashKey)
	h.Write([]byte(id))
	dst := h.Sum(nil)

	return &http.Cookie{
		Name:   "token",
		Value:  fmt.Sprintf("%s:%x", id, dst),
		Path:   "/",
		MaxAge: (int)(time.Hour * 24 * 30),
	}
}

func GetID(r *http.Request) (string, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return "", err
	}
	parts := strings.Split(cookie.Value, ":")

	if len(parts) != 2 {
		return "", errors.New("wrong token format")
	}
	id := parts[0]
	key := parts[1]

	h := hmac.New(sha256.New, hashKey)
	h.Write([]byte(id))
	dst := h.Sum(nil)

	if !hmac.Equal([]byte(key), dst) {
		return "", errors.New("wrong token")
	}
	return id, nil
}

func Check(r *http.Request) bool {
	_, err := GetID(r)
	return err != nil
}

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

func init() {
	key, err := generateRandom(16)
	if err != nil {
		panic(err)
	}
	hashKey = key
}
