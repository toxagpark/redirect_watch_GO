package main

import (
	"fmt"
	"log"
	"net/http"
	"url_shortener_go/handler"
	"url_shortener_go/postgres"
	"url_shortener_go/redis"

	"github.com/gorilla/mux"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	postgres.InitDB()
	redis.InitRedis()

	r := mux.NewRouter()
	r.HandleFunc("/shorten", handler.ShorterHandler).Methods("POST")
	r.HandleFunc("/{short_code}", handler.RedirectHandler).Methods("GET")

	r.Use(loggingMiddleware)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
