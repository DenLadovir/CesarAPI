package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Task struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func main() {
	task := Task{
		Title:       "Ещё одна задача",
		Description: "Описание для второй тестовой задачи",
		Completed:   false,
	}

	taskJSON, err := json.Marshal(task)
	if err != nil {
		fmt.Println("Ошибка при сериализации задачи:", err)
		return
	}

	resp, err := http.Post("http://localhost:8000/tasks", "application/json", bytes.NewBuffer(taskJSON))
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated || resp.Status == "200 OK" {
		fmt.Println("Задача успешно добавлена!")
	} else {
		fmt.Printf("Ошибка при добавлении задачи: %s\n", resp.Status)
	}
}
