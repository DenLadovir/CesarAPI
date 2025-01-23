package telegram

import (
	"bytes"
	"fmt"
	"net/http"
)

const telegramAPI = "https://api.telegram.org/bot"

// SendTelegramMessage отправляет сообщение в Telegram
func SendTelegramMessage(botToken string, chatID int64, message string) error {
	url := fmt.Sprintf("%s%s/sendMessage", telegramAPI, botToken)

	// Формируем тело запроса
	body := []byte(fmt.Sprintf(`{"chat_id":%d,"text":"%s"}`, chatID, message))

	// Создаем HTTP-запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}

	return nil
}
