package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Структура ошибки
type ApiError struct {
	Error string `json:"error"`
}

// Структура респонса
type GoodResponse struct {
	Status string `json:"status"`
}

// Структура заметки
type Note struct {
	Id           int    `json:"id"`
	UserId       int    `json:"userId"`
	Note         string `json:"note"`
	CreationTime string `json:"creationTime"`
}

// Функция создания заметки
func createNote(db *sqlx.DB, values url.Values) error {

	// Проверка на наличие параметров
	if !values.Has("userId") || !values.Has("note") {
		return errors.New("userId or note not found")
	}

	// Добавление заметки в базу
	_, err := db.Exec("INSERT INTO notes (userId, note, creationTime) values ($1, $2, $3)", values.Get("userId"), values.Get("note"), time.Now().Format("2006-01-02T15:04:05"))
	if err != nil {
		return err
	}

	return nil
}

// Роут "/create"
func Create(w http.ResponseWriter, r *http.Request) {

	// Передача в заголовок респонса типа данных
	w.Header().Set("Content-Type", "application/json")

	// Проверка на попытку SQL-инъекций
	if strings.ContainsAny(r.URL.String(), "%'") {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(ApiError{Error: "Bad Request"})
		w.Write(json)
		return
	}

	// Подключение к БД
	db, err := ConnectDB()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(ApiError{Error: "Internal Server Error"})
		w.Write(json)
		log.Printf("connectDB error: %s", err)
		return
	}

	// Создание заметки
	err = createNote(db, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(ApiError{Error: "Internal Server Error"})
		w.Write(json)
		log.Printf("createNote error: %s", err)
		return
	}

	// Получение статистики, форматирование и отправка
	jsonResp, err := json.Marshal(GoodResponse{Status: "Created"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(ApiError{Error: "Internal Server Error"})
		w.Write(json)
		log.Printf("json.Marshal error: %s", err)
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResp)
	}

}
