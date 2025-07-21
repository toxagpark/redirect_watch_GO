package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"url_shortener_go/postgres"
	"url_shortener_go/redis"
)

func ShorterHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON format"))
		return
	}

	short_code := postgres.AddURLPostgres(input.URL)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(short_code))
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	short_сode := vars["short_code"]
	//ищем в редисе
	log.Println("ищем в редисе")
	original_url := redis.GetURLRedis(short_сode)
	//если не нашли ищем в постгресе
	if original_url == "" {
		log.Println("Ищем в постгресе")
		original_url = postgres.GetURLPostgres(short_сode)
	}
	//если совсем не нашли - ошибка
	if original_url == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("nothing in database"))
	}
	//добавляем в редис
	redis.AddURLRedis(original_url, short_сode)
	//добавляем гостя
	ip := r.RemoteAddr
	userAgent := r.UserAgent()
	postgres.AddVisitPostgres(ip, userAgent, original_url)
	http.Redirect(w, r, original_url, http.StatusMovedPermanently)
}
