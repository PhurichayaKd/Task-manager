package service

import (
	"context"
	"task-manager/internal/auth"
	"task-manager/internal/domain"
	"task-manager/internal/repo"
)

type UserService interface {
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	UpdateName(ctx context.Context, id int64, name string) error
	UpdateUsername(ctx context.Context, id int64, username string) error
	UpdatePassword(ctx context.Context, id int64, password string) error
}

type userService struct {
	userRepo repo.UserRepo
	hasher   auth.PasswordHasher
}

func NewUserService(userRepo repo.UserRepo, hasher auth.PasswordHasher) UserService {
	return &userService{
		userRepo: userRepo,
		hasher:   hasher,
	}
}

func (s *userService) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) UpdateName(ctx context.Context, id int64, name string) error {
	return s.userRepo.UpdateName(ctx, id, name)
}

func (s *userService) UpdateUsername(ctx context.Context, id int64, username string) error {
	return s.userRepo.UpdateUsername(ctx, id, username)
}

func (s *userService) UpdatePassword(ctx context.Context, id int64, password string) error {
	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		return err
	}
	return s.userRepo.UpdatePassword(ctx, id, hashedPassword)
}
