package repo

import (
	"context"
	"database/sql"

	"task-manager/internal/domain"
)

type TaskRepo interface {
	GetByUserID(ctx context.Context, userID int, limit int) ([]*domain.Task, error)
	Create(ctx context.Context, task *domain.Task) (*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, id int, userID int) error
}

type taskRepo struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) TaskRepo {
	return &taskRepo{db: db}
}

func (r *taskRepo) GetByUserID(ctx context.Context, userID int, limit int) ([]*domain.Task, error) {
	// For now, return empty slice to prevent errors
	return []*domain.Task{}, nil
}

func (r *taskRepo) Create(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	// Placeholder implementation
	return task, nil
}

func (r *taskRepo) Update(ctx context.Context, task *domain.Task) error {
	// Placeholder implementation
	return nil
}

func (r *taskRepo) Delete(ctx context.Context, id int, userID int) error {
	// Placeholder implementation
	return nil
}
