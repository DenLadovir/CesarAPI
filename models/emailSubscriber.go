package models

import (
	"gorm.io/gorm"
)

// EmailSubscriber представляет подписчика электронной почты
type EmailSubscriber struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Email string `json:"email" gorm:"unique;not null"`
}

// AddEmailSubscriber добавляет новый адрес электронной почты в базу данных
func AddEmailSubscriber(db *gorm.DB, email string) error {
	subscriber := EmailSubscriber{Email: email}
	return db.Create(&subscriber).Error
}

// GetEmailSubscribers извлекает все адреса электронной почты из базы данных
func GetEmailSubscribers(db *gorm.DB) ([]string, error) {
	var subscribers []EmailSubscriber
	if err := db.Find(&subscribers).Error; err != nil {
		return nil, err
	}

	var emails []string
	for _, subscriber := range subscribers {
		emails = append(emails, subscriber.Email)
	}
	return emails, nil
}
