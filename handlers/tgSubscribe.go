package handlers

import (
	"CesarAPI/models"
	"CesarAPI/telegram"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

func GetTelegramChannels(db *gorm.DB) ([]models.TelegramChannel, error) {
	var channels []models.TelegramChannel
	if err := db.Find(&channels).Error; err != nil {
		return nil, err
	}
	return channels, nil
}

func AddTelegramChannel(db *gorm.DB, token string, chatID int64) error {
	channel := models.TelegramChannel{Token: token, ChatID: chatID}
	return db.Create(&channel).Error
}

func NotifyTelegramChannels(db *gorm.DB, message string) {
	var channels []models.TelegramChannel

	// Извлекаем все Telegram-каналы из базы данных
	if err := db.Find(&channels).Error; err != nil {
		log.Printf("Не удалось получить список Telegram-каналов: %v\n", err)
		return
	}

	// Отправляем сообщение каждому каналу
	for _, channel := range channels {
		err := telegram.SendTelegramMessage(channel.Token, channel.ChatID, message)
		if err != nil {
			log.Printf("Не удалось отправить сообщение в Telegram канал %d: %v\n", channel.ChatID, err)
		} else {
			log.Printf("Сообщение успешно отправлено в Telegram канал %d", channel.ChatID)
		}
	}
}

func AddTelegramChannelHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request models.TelegramChannelRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		channel := models.TelegramChannel{
			Token:  request.Token,
			ChatID: request.ChatID,
		}

		// Сохранение в базу данных
		if err := db.Create(&channel).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Channel added successfully"})
	}
}
