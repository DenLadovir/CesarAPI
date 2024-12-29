package models

const (
	StatusPending    = "Pending"
	StatusInProgress = "In_progress"
	StatusCompleted  = "Completed"
)

type Task struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Version     int    `json:"version"`
}