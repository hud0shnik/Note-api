package main

import (
	"fmt"
	"log"
	"net/http"
	"note-api/api"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// Функция для ручного запуска на 127.0.0.1
func main() {

	// Вывод времени начала работы
	fmt.Println("API Start: " + string(time.Now().Format("2006-01-02 15:04:05")))
	fmt.Println("Port:\t" + os.Getenv("PORT"))

	// Роутер
	router := mux.NewRouter()

	// Маршруты
	router.HandleFunc("/api/search", api.Search).Methods("GET")
	router.HandleFunc("/api/create", api.Search).Methods("POST")

	// Запуск API
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))

}
