package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
)

type Links map[string]string

var links Links

func main() {
	links = make(Links)
	http.HandleFunc(`/`, handleRoot)
	if err := http.ListenAndServe(`:8080`, nil); err != nil {
		panic(err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handlePost(w, r)
	}
	if r.Method == http.MethodGet {
		handleGet(w, r)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	body := readBody(r)
	fmt.Println(body)

	var key string
	for {
		key = generateKey()
		if _, exists := links[key]; !exists {
			break
		}
	}

	links[key] = body
	parts := []string{`http:/`, r.Host, key}
	result := strings.Join(parts, `/`)

	w.WriteHeader(201)
	_, err := w.Write([]byte(result))
	if err != nil {
		return
	}
	fmt.Println(key, body)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), `/`)
	if len(parts) != 2 {
		w.WriteHeader(404)
		return
	}
	key := parts[1]
	url, ok := links[key]
	if !ok {
		w.WriteHeader(404)
		return
	}

	w.Header().Set(`Location`, url)
	w.WriteHeader(307)

	fmt.Println(url)
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

func generateKey() string {
	return generateRandomString(8)
}

func generateRandomString(length int) string {
	charset := `ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789`
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
