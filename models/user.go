package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"column:username;unique;not null"`
	Password string `gorm:"not null"`
}
