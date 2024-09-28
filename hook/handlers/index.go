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

	headers := http.Header{}
	for key := range response.Headers {
		if !(strings.HasPrefix(key, "Cf-") || key == "Cdn-Loop" || key == "X-Real-IP" || key == "X-Forwarded-For" || key == "X-Forwarded-Proto") {
			headers.Add(key, response.Headers[key])
		}
	}

	db.PublishRequest(bin, db.Request{
		Method:  r.Method,
		URL:     r.URL.String(),
		Path:    r.URL.Path,
		Sender:  r.Header.Get("X-Real-IP"),
		Query:   r.URL.Query(),
		Headers: headers,
		Body:    string(body),
		Time:    time.Now(),
	})

	for key, value := range response.Headers {
		w.Header().Add(key, value)
	}
	w.WriteHeader(response.StatusCode)
	w.Write([]byte(response.Body))
}
