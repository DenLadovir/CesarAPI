package email

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

// SendEmail отправляет электронное письмо
func SendEmail(emailList []string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "denis.ladovir@yandex.ru")
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.yandex.ru", 465, "denis.ladovir@yandex.ru", "pevkuznujhzmdjdp")

	for _, email := range emailList {
		m.SetHeader("To", email) // Устанавливаем адрес получателя
		log.Printf("Отправка письма на %s с темой: %s", email, subject)

		// Отправка письма
		if err := d.DialAndSend(m); err != nil {
			log.Println("Ошибка при отправке письма:", err)
			return err
		}

		log.Printf("Письмо успешно отправлено на %s", email)
	}

	return nil
}
