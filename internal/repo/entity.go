package repo

import "time"

// Task - структура, соответствующая таблице tasks
type Task struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GetTask struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ID          int       `json:"id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
