package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Функция поиска заметок
func SearchNotes(db *sqlx.DB, values url.Values) ([]Note, error) {

	// Начало запроса и слайс параметров
	query := "SELECT * FROM notes"
	parameters := []string{}

	// Проверки на наличие параметров и запись их в слайс
	if values.Has("userId") {
		parameters = append(parameters, "userId = "+values.Get("userId"))
	} else {
		return nil, errors.New("userId not found")
	}
	if values.Has("id") {
		parameters = append(parameters, "id = "+values.Get("id"))
	}

	query += " WHERE " + strings.Join(parameters, " AND ")

	// Инициализация результата
	var result []Note

	// Получение и проверка данных
	err := db.Select(&result, query+" ORDER BY id DESC")
	if err != nil {
		return nil, err
	} else if len(result) == 0 {
		return nil, errors.New("notes not found")
	}

	return result, nil

}

// Функция подключения к БД
func ConnectDB() (*sqlx.DB, error) {

	// Инициализация переменных окружения
	godotenv.Load()

	// Подключение к БД
	db, err := sqlx.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USERNAME"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PASSWORD")))
	if err != nil {
		return nil, errors.New("failed to connect")
	}

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		return nil, errors.New("failed to ping db")
	}

	return db, nil
}

// Роут "/search"
func Search(w http.ResponseWriter, r *http.Request) {

	// Передача в заголовок респонса типа данных
	w.Header().Set("Content-Type", "application/json")

	// Подключение к БД
	db, err := ConnectDB()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(ApiError{Error: "Internal Server Error"})
		w.Write(json)
		log.Printf("connectDB error: %s", err)
		return
	}

	// Поиск заметок
	notes, err := SearchNotes(db, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(ApiError{Error: "Internal Server Error"})
		w.Write(json)
		log.Printf("searchNotes error: %s", err)
		return
	}

	// Форматирование и отправка заметок
	jsonResp, err := json.Marshal(notes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(ApiError{Error: "Internal Server Error"})
		w.Write(json)
		log.Printf("json.Marshal error: %s", err)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	}

}
