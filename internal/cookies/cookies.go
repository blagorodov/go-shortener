package cookies

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

var hashKey = []byte("secrethashkey")

type key string

const ContextKey key = "userID"

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

func GetID(w http.ResponseWriter, r *http.Request) (string, error) {
	var token string
	c, err := r.Cookie("token")
	if err != nil {
		token = r.Header.Get("Authorization")
	} else {
		token = c.Value
	}
	if token == "" {
		cookie := New()
		token = cookie.Value
		http.SetCookie(w, cookie)
		w.Header().Set("Authorization", token)
	}
	return GetIDToken(token)
}

func GetIDToken(token string) (string, error) {
	parts := strings.Split(token, ":")

	if len(parts) != 2 {
		return "", errors.New("wrong token format")
	}
	id := parts[0]
	key, _ := hex.DecodeString(parts[1])

	h := hmac.New(sha256.New, hashKey)
	h.Write([]byte(id))
	dst := h.Sum(nil)

	if !hmac.Equal(dst, key) {
		return "", errors.New("wrong token")
	}
	return id, nil
}

//func Check(r *http.Request) bool {
//	_, err := GetID(r)
//	return err == nil
//}

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
