package service

import (
	"context"

	"task-manager/internal/domain"
	"task-manager/internal/repo"
)

type TaskService interface {
	GetUserTasks(ctx context.Context, userID int, limit int) ([]*domain.Task, error)
	CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error)
	UpdateTask(ctx context.Context, task *domain.Task) error
	DeleteTask(ctx context.Context, id int, userID int) error
}

type taskService struct {
	taskRepo repo.TaskRepo
}

func NewTaskService(taskRepo repo.TaskRepo) TaskService {
	return &taskService{taskRepo: taskRepo}
}

func (s *taskService) GetUserTasks(ctx context.Context, userID int, limit int) ([]*domain.Task, error) {
	return s.taskRepo.GetByUserID(ctx, userID, limit)
}

func (s *taskService) CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	return s.taskRepo.Create(ctx, task)
}

func (s *taskService) UpdateTask(ctx context.Context, task *domain.Task) error {
	return s.taskRepo.Update(ctx, task)
}

func (s *taskService) DeleteTask(ctx context.Context, id int, userID int) error {
	return s.taskRepo.Delete(ctx, id, userID)
}
