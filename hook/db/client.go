package db

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	BinPrefix     = "bin:"
	ChannelPrefix = "channel:"

	DefaultHTTPStatusCode = "200"
	DefaultHTTPBody       = "Hello, World!"
	DefaultHTTPHeaders    = `{"Content-Type": "text/plain"}`
)

var (
	Client *redis.Client
	Ctx    = context.Background()
)

type Request struct {
	Method  string
	URL     string
	Path    string
	Sender  string
	Query   map[string][]string
	Headers map[string][]string
	Body    string
	Time    time.Time
}

// User controllable HTTP response for a bin
type Response struct {
	Body       string
	StatusCode int
	Headers    map[string]string
}

func InitRedis(url string) {
	Client = redis.NewClient(&redis.Options{
		Addr: url,
	})
}

func AllocateDefaultBin() (string, error) {
	bin, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyz", 14)
	if err != nil {
		return "", err
	}

	binHashKey := BinPrefix + bin
	err = Client.HSet(Ctx, binHashKey, []string{
		"body", DefaultHTTPBody,
		"status_code", DefaultHTTPStatusCode,
		"headers", DefaultHTTPHeaders,
	}).Err()
	if err != nil {
		return "", err
	}
	return bin, Client.Expire(Ctx, binHashKey, time.Duration(24)*time.Hour).Err()
}

func GetResponseForBin(bin string) (Response, error) {
	binHashKey := BinPrefix + bin
	response, err := Client.HGetAll(Ctx, binHashKey).Result()
	if err != nil {
		return Response{}, err
	}

	body := response["body"]
	statusCode, err := strconv.Atoi(response["status_code"])
	if err != nil {
		return Response{}, err
	}
	headers := make(map[string]string)
	if err := json.Unmarshal([]byte(response["headers"]), &headers); err != nil {
		return Response{}, err
	}

	return Response{
		Body:       body,
		StatusCode: statusCode,
		Headers:    headers,
	}, nil
}

func PublishRequest(bin string, request Request) error {
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return err
	}

	return Client.Publish(Ctx, ChannelPrefix+bin, requestJSON).Err()
}
