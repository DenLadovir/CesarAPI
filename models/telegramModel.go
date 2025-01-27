package models

type TelegramChannel struct {
	ID     uint   `gorm:"primaryKey"`
	Token  string `gorm:"not null"`
	ChatID int64  `gorm:"not null"`
}

type TelegramChannelRequest struct {
	Token  string `json:"token"`
	ChatID int64  `json:"chat_id"`
}
