package models

const (
	StatusPending    = "Pending"
	StatusInProgress = "In_progress"
	StatusCompleted  = "Completed"
)

type Task struct {
	ID           int    `json:"id" gorm:"primaryKey"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Status       string `json:"status"`
	Version      int    `json:"version"`
	UpdateByUser string `json: "update_by_user"`
	UpdateTime   string `json: "update_time"`
}

type UpdateTask struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
	Version     int     `json:"version"`
}
