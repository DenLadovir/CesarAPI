package routes

import (
	"CesarAPI/database"
	"CesarAPI/email"
	"CesarAPI/models"
	"CesarAPI/telegram"
	"CesarAPI/utils"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type TaskHandler struct{}

func (h *TaskHandler) RegisterRoutes(r *chi.Mux, db *gorm.DB) {
	r.Get("/tasks", h.GetTasks)
	r.Post("/tasks", h.CreateTask)
	r.Put("/tasks/{id}", h.UpdateTask)
	r.Delete("/tasks/{id}", h.DeleteTask)
	r.Put("/tasks/{id}/status", func(w http.ResponseWriter, r *http.Request) {
		h.UpdateTaskStatus(w, r, db) // Передаем db в метод
	})
	r.Post("/register", h.RegisterHandler(db))
	r.Post("/login", h.LoginHandler(db))
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	database.DB.Find(&tasks)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("Error decoding task: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.DB.Create(&task).Error; err != nil {
		log.Printf("Error creating task: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	emailBody := "Новая задача создана:\n\n" +
		"Название: " + task.Title + "\n" +
		"Описание: " + task.Description + "\n" +
		"Статус: " + task.Status
	emailList := []string{"denis.ladovir@yandex.ru", "denis120992@gmail.com", "den-12.09.92@mail.ru"}
	err := email.SendEmail(emailList, "Новая задача создана", emailBody)
	if err != nil {
		log.Printf("Не удалось отправить уведомление по электронной почте: %v\n", err)
	}

	botToken := ""             // Замените на ваш токен телеграм бота
	chatID := int64(123456789) // Замените на ваш chat_id
	telegramMessage := fmt.Sprintf("Новая задача создана:\nНазвание: %s\nОписание: %s\nСтатус: %s", task.Title, task.Description, task.Status)
	err = telegram.SendTelegramMessage(botToken, chatID, telegramMessage)
	if err != nil {
		log.Printf("Не удалось отправить уведомление в Telegram: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") // Получаем ID из URL
	var updateTask models.UpdateTask

	// Декодируем тело запроса в структуру UpdateTask
	if err := json.NewDecoder(r.Body).Decode(&updateTask); err != nil {
		log.Printf("Updating error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Находим существующую задачу по ID
	var existingTask models.Task
	if err := database.DB.First(&existingTask, id).Error; err != nil {
		log.Printf("Task not found: %v", err)
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Проверяем версию
	if existingTask.Version != updateTask.Version {
		log.Printf("Version conflict: existing version %d, provided version %d", existingTask.Version, updateTask.Version)
		http.Error(w, "Version conflict: task was modified by another user", http.StatusConflict)
		return
	}

	// Обновляем поля только если они указаны
	if updateTask.Title != nil {
		existingTask.Title = *updateTask.Title
	}
	if updateTask.Description != nil {
		existingTask.Description = *updateTask.Description
	}
	if updateTask.Status != nil {
		existingTask.Status = *updateTask.Status
	}

	// Увеличиваем версию
	existingTask.Version++
	existingTask.UpdateByUser = "Тестовый пользователь 34"
	existingTask.UpdateTime = time.Now().Format("02.01.2006 15:04:05")

	if err := database.DB.Save(&existingTask).Error; err != nil {
		log.Printf("Updating error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Добавляем рассылку на электронные почты при изменении задачи
	emailBody := "Задача обновлена:\n\n" +
		"Название: " + existingTask.Title + "\n" +
		"Описание: " + existingTask.Description + "\n" +
		"Статус: " + existingTask.Status
	emailList := []string{"denis.ladovir@yandex.ru", "denis120992@gmail.com", "den-12.09.92@mail.ru"}
	err := email.SendEmail(emailList, "Задача обновлена", emailBody)
	if err != nil {
		log.Printf("Не удалось отправить уведомление по электронной почте: %v\n", err)
	}

	botToken := "6265294268:AAHR_FepTPa_3u8rTbRU9ke6ELLMMPeh4RA" // Замените на ваш токен телеграм бота
	chatID := int64(476899260)                                   // Замените на ваш chat_id https://api.telegram.org/bot<Ваш_Токен>/getUpdates
	telegramMessage := fmt.Sprintf("Задача обновлена:\nНазвание: %s\nОписание: %s\nСтатус: %s", existingTask.Title, existingTask.Description, existingTask.Status)
	err = telegram.SendTelegramMessage(botToken, chatID, telegramMessage)
	if err != nil {
		log.Printf("Не удалось отправить уведомление в Telegram: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingTask)

	//authHeader := r.Header.Get("Authorization")
	//if authHeader == "" {
	//	http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
	//	return
	//}
	//
	//tokenString := strings.Split(authHeader, " ")[1]
	//claims := &models.Claims{}
	//token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
	//	return []byte("your_secret_key"), nil
	//})
	//
	//if err != nil || !token.Valid {
	//	http.Error(w, "Invalid token", http.StatusUnauthorized)
	//	return
	//}
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Ошибка преобразования ID: %v\n", err)
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	log.Printf("Попытка удалить задачу с ID: %d\n", id)

	// Удаляем задачу по ID
	if err := database.DB.Delete(&models.Task{}, id).Error; err != nil {
		log.Printf("Delete error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Задача с ID %d была успешно удалена.\n", id)

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) UpdateTaskStatus(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Ошибка преобразования ID: %v\n", err)
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var input struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("Ошибка декодирования JSON: %v\n", err)
		http.Error(w, "Неверный ввод", http.StatusBadRequest)
		return
	}

	// Находим задачу по ID
	var task models.Task
	if err := db.First(&task, id).Error; err != nil {
		log.Printf("Задача не найдена: %v\n", err)
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	// Обновляем статус задачи
	task.Status = input.Status
	if err := db.Save(&task).Error; err != nil {
		log.Printf("Ошибка обновления статуса: %v\n", err)
		http.Error(w, "Не удалось обновить статус задачи", http.StatusInternalServerError)
		return
	}

	log.Printf("Статус задачи с ID %d был успешно обновлён на %s.\n", id, input.Status)
	w.WriteHeader(http.StatusNoContent) // Успешное обновление
}

func (h *TaskHandler) RegisterHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User // Предполагается, что у вас есть структура User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		var existingUser models.User
		if err := db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
			http.Error(w, "User already exists!", http.StatusConflict)
			return
		}

		// Хеширование пароля
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}
		user.Password = string(hashedPassword)

		// Сохраните пользователя в базе данных
		result := db.Create(&user) // Используем db.Create для автоматического заполнения полей
		if result.Error != nil {
			http.Error(w, "Error saving user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User registered successfully!")) // Можно добавить сообщение об успешной регистрации
	}
}

func (h *TaskHandler) LoginHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		// Декодирование JSON-данных из запроса
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		var user models.User
		// Поиск пользователя в базе данных
		result := db.Where("username = ?", credentials.Username).First(&user)
		if result.Error != nil {
			http.Error(w, "Invalid credentials!", http.StatusUnauthorized)
			return
		}

		// Проверка пароля
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
		if err != nil {
			http.Error(w, "Invalid credentials!", http.StatusUnauthorized)
			return
		}

		// Генерация JWT
		token, err := utils.GenerateJWT(user.Username)
		if err != nil {
			http.Error(w, "Could not generate token", http.StatusInternalServerError)
			return
		}

		// Успешная аутентификация
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token)) // Возвращаем токен клиенту
	}
}
