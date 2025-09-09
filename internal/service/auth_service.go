package service

import (
	"context"
	"database/sql"
	"strings"

	"task-manager/internal/auth"
	"task-manager/internal/domain"
	"task-manager/internal/repo"
)

type AuthService interface {
	// คืน 4 ค่า: user, accessToken, created(สมัครใหม่?), error
	LoginOrSignupGoogle(ctx context.Context, email, name, sub, avatar string) (*domain.User, string, bool, error)
	Login(ctx context.Context, usernameOrEmail, password string) (*domain.User, error)
	Register(ctx context.Context, email, username, password, name string) (*domain.User, error)
	CompleteGoogleRegistration(ctx context.Context, email, username, password, name string) (*domain.User, error)
	GenerateToken(userID int64) (string, error)
}

type authService struct {
	UserRepo repo.UserRepo
	Hasher   auth.PasswordHasher
	JWT      auth.JWT
}

func NewAuthService(ur repo.UserRepo, hasher auth.PasswordHasher, jwt auth.JWT) AuthService {
	return &authService{UserRepo: ur, Hasher: hasher, JWT: jwt}
}

func ns(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func (s *authService) LoginOrSignupGoogle(ctx context.Context, email, name, sub, avatar string) (*domain.User, string, bool, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	// มีอยู่แล้ว?
	u, err := s.UserRepo.GetByEmail(ctx, email)
	if err != nil && err != repo.ErrNotFound {
		return nil, "", false, err
	}

	created := false
	if u == nil {
		// สมัครใหม่
		u = &domain.User{
			Email:      email,
			Name:       ns(name),
			Provider:   ns("google"),
			ProviderID: ns(sub),
			AvatarURL:  ns(avatar),
			Role:       "user",
		}
		id, err := s.UserRepo.CreateFromOAuth(ctx, u)
		if err != nil {
			return nil, "", false, err
		}
		u.ID = id
		created = true
	}

	// ออก access token
	token, _, err := s.JWT.GenerateAccessToken(u.ID, u.Role) // ฟังก์ชันนี้คืน (token, ttl, err)
	if err != nil {
		return nil, "", created, err
	}

	return u, token, created, nil
}

func (s *authService) Login(ctx context.Context, usernameOrEmail, password string) (*domain.User, error) {
	usernameOrEmail = strings.TrimSpace(usernameOrEmail)

	if usernameOrEmail == "" || password == "" {
		return nil, domain.ErrInvalidInput
	}

	// Try to get user by email first, then by username
	var user *domain.User
	var err error

	// Check if input looks like an email (contains @)
	if strings.Contains(usernameOrEmail, "@") {
		user, err = s.UserRepo.GetByEmail(ctx, strings.ToLower(usernameOrEmail))
	} else {
		user, err = s.UserRepo.GetByUsername(ctx, usernameOrEmail)
	}

	if err != nil {
		if err == repo.ErrNotFound {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	// Check password
	if !user.PasswordHash.Valid {
		return nil, domain.ErrInvalidCredentials
	}

	if err := s.Hasher.Compare(user.PasswordHash.String, password); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	return user, nil
}

func (s *authService) GenerateToken(userID int64) (string, error) {
	token, _, err := s.JWT.GenerateAccessToken(userID, "user")
	return token, err
}

func (s *authService) Register(ctx context.Context, email, username, password, name string) (*domain.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	username = strings.TrimSpace(username)
	name = strings.TrimSpace(name)

	if email == "" || username == "" || password == "" || name == "" {
		return nil, domain.ErrInvalidInput
	}

	// Check if email or username already exists
	emailExists, err := s.UserRepo.EmailExists(ctx, email)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, domain.ErrEmailAlreadyExists
	}

	usernameExists, err := s.UserRepo.UsernameExists(ctx, username)
	if err != nil {
		return nil, err
	}
	if usernameExists {
		return nil, domain.ErrUsernameAlreadyExists
	}

	// Hash password
	hashedPassword, err := s.Hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
		Email:        email,
		Username:     ns(username),
		PasswordHash: ns(hashedPassword),
		Name:         ns(name),
		Role:         "user",
	}

	createdUser, err := s.UserRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *authService) CompleteGoogleRegistration(ctx context.Context, email, username, password, name string) (*domain.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	username = strings.TrimSpace(username)
	name = strings.TrimSpace(name)

	if email == "" || username == "" || password == "" || name == "" {
		return nil, domain.ErrInvalidInput
	}

	// Check if email already exists
	existingUser, err := s.UserRepo.GetByEmail(ctx, email)
	if err != nil && err != repo.ErrNotFound {
		return nil, err
	}
	
	// If user exists, update their information
	if existingUser != nil {
		// Check if username already exists (but not for this user)
		usernameExists, err := s.UserRepo.UsernameExists(ctx, username)
		if err != nil {
			return nil, err
		}
		if usernameExists {
			return nil, domain.ErrUsernameAlreadyExists
		}

		// Hash password
		hashedPassword, err := s.Hasher.Hash(password)
		if err != nil {
			return nil, err
		}

		// Update user with username and password
		existingUser.Username = ns(username)
		existingUser.PasswordHash = ns(hashedPassword)
		existingUser.Name = ns(name)

		updatedUser, err := s.UserRepo.Update(ctx, existingUser)
		if err != nil {
			return nil, err
		}

		return updatedUser, nil
	}
	
	// If user doesn't exist, create a new one
	return s.Register(ctx, email, username, password, name)
}
