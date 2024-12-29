package routes

import (
	"CesarAPI/database"
	"CesarAPI/models"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") // Получаем ID из URL
	var task models.Task

	// Декодируем тело запроса в структуру Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
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
	if existingTask.Version != task.Version {
		log.Printf("Version conflict: existing version %d, provided version %d", existingTask.Version, task.Version)
		http.Error(w, "Version conflict: task was modified by another user", http.StatusConflict)
		return
	}

	// Обновляем задачу
	task.Version++            // Увеличиваем версию
	task.ID = existingTask.ID // Устанавливаем ID для обновления
	if err := database.DB.Save(&task).Error; err != nil {
		log.Printf("Updating error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
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