package handlers

import (
	"io"
	"net/http"
	"strings"
	"time"
	"webhooks/db"
)

func Index(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	parts := strings.Split(host, ".")
	if len(parts) < 4 {
		http.Error(w, "Invalid host format", http.StatusBadRequest)
		return
	}

	bin := parts[0]
	response, err := db.GetResponseForBin(bin)
	if err != nil {
		http.Error(w, "Bin not found", http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		body = []byte{}
	}

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
		Sender:  r.Header.Get("X-Real-IP"),
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
