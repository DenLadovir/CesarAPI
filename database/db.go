package database

import (
	"CesarAPI/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() error {
	var err error
	dsn := "host=localhost user=D10n password=123789456 dbname=tasks port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Connection error: %v", err)
		return err
	}
	DB.AutoMigrate(&models.Task{})
	log.Println("Успешно подключено к базе данных.")
	return nil
}
