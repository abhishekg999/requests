package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"webhooks/db"
	"webhooks/handlers"
)

func main() {
	db.InitRedis(os.Getenv("REDIS_URL"))
	http.HandleFunc("/", handlers.Index)

	fmt.Println("Starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
