package handlers

import (
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"webhooks/db"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func FullURL(r *http.Request) string {
	return r.Host + r.URL.String()
}

var BIN_REGEX = regexp.MustCompile(GetEnv("BIN_REGEX", `^(?:[a-zA-Z0-9-_.]+?)(?:\:\d+)?/(?<bin>\w*)(?:/.*)?$`))

func Index(w http.ResponseWriter, r *http.Request) {
	url := FullURL(r)
	match := BIN_REGEX.FindStringSubmatch(url)

	if match == nil || len(match) != 2 {
		http.Error(w, "Invalid url format", http.StatusBadRequest)
		return
	}
	bin := match[1]
	response, err := db.GetResponseForBin(bin)
	if err != nil {
		http.Error(w, "Bin: "+bin+" not found", http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		body = []byte{}
	}

	Sender := r.Header.Get("X-Real-IP")

	for key := range r.Header {
		key_lower := strings.ToLower(key)
		if strings.HasPrefix(key_lower, "cf-") || key_lower == "cdn-loop" || key_lower == "x-real-ip" || key_lower == "x-forwarded-for" || key_lower == "x-forwarded-proto" {
			r.Header.Del(key)
		}
	}

	db.PublishRequest(bin, db.Request{
		Method:  r.Method,
		URL:     r.URL.String(),
		Path:    r.URL.Path,
		Sender:  Sender,
		Query:   r.URL.Query(),
		Headers: r.Header,
		Body:    string(body),
		Time:    time.Now(),
	})

	for key, value := range response.Headers {
		w.Header().Add(key, value)
	}
	w.WriteHeader(response.StatusCode)
	w.Write([]byte(response.Body))
}
