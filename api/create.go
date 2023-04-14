package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
)

// Структура заметки
type Note struct {
	Id     int    `json:"id"`
	UserId int    `json:"userId"`
	Note   string `json:"note"`
	Time   string `json:"time"`
}

// Функция создания заметки
func CreateNote(db *sqlx.DB, values url.Values) error {

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

	// Подключение к БД
	db, err := ConnectDB()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(map[string]string{"Error": err.Error()})
		w.Write(json)
		log.Printf("connectDB error: %s", err)
	}

	// Получение статистики, форматирование и отправка
	jsonResp, err := json.Marshal(CreateNote(db, r.URL.Query()))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(map[string]string{"Error": "Internal Server Error"})
		w.Write(json)
		log.Printf("json.Marshal error: %s", err)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	}

}
