package handlers

import (
	"CesarAPI/models"
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
)

// SubscribeHandler представляет обработчик подписки
type SubscribeHandler struct {
	DB *gorm.DB // Инициализация базы данных
}

// NewSubscribeHandler создает новый SubscribeHandler
func NewSubscribeHandler(db *gorm.DB) *SubscribeHandler {
	return &SubscribeHandler{DB: db}
}

// HandleSubscribe обрабатывает запросы на подписку
func (h *SubscribeHandler) HandleSubscribe(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Email string `json:"email"`
	}

	// Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Добавляем адрес электронной почты в базу данных
	if err := models.AddEmailSubscriber(h.DB, requestBody.Email); err != nil {
		http.Error(w, "Could not add subscriber", http.StatusInternalServerError)
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Subscriber added"))
}
