package email

import (
	"gopkg.in/gomail.v2"
	"log"
)

// SendEmail отправляет электронное письмо
func SendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "denis.ladovir@yandex.ru") // Укажите ваш адрес электронной почты
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.yandex.ru", 465, "denis.ladovir@yandex.ru", "pevkuznujhzmdjdp")
	log.Printf("Отправка письма на %s с темой: %s", to, subject)

	// Отправка письма
	if err := d.DialAndSend(m); err != nil {
		log.Println("Ошибка при отправке письма:", err)
		return err
	}
	return nil
}
