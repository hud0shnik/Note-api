package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/jmoiron/sqlx"
)

// Функция удаления заметки
func deleteNote(db *sqlx.DB, values url.Values) error {
	return nil
}

// Роут "/delete"
func Delete(w http.ResponseWriter, r *http.Request) {

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

	// Удаление заметки
	err = deleteNote(db, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(ApiError{Error: "Internal Server Error"})
		w.Write(json)
		log.Printf("searchNotes error: %s", err)
		return
	}

	// Форматирование и отправка заметок
	jsonResp, err := json.Marshal(GoodResponse{Status: "Successfully deleted"})
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
