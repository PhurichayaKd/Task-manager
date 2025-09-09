package domain

import (
	"database/sql"
	"time"
)

// Task represents a task in the system
type Task struct {
	ID          int            `json:"id" db:"id"`
	UserID      int            `json:"user_id" db:"user_id"`
	Title       string         `json:"title" db:"title"`
	Description sql.NullString `json:"description" db:"description"`
	Status      string         `json:"status" db:"status"` // pending, in_progress, completed
	Priority    string         `json:"priority" db:"priority"` // low, medium, high
	DueDate     sql.NullTime   `json:"due_date" db:"due_date"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}
