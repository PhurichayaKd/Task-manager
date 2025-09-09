package auth

import (
	"errors"
	"time"

	"task-manager/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int64  `json:"uid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWT interface {
	GenerateAccessToken(userID int64, role string) (string, time.Duration, error)
	GenerateRefreshToken(userID int64) (string, time.Duration, error)
	ParseAccess(token string) (*Claims, error)
	ValidateToken(token string) (*Claims, error)
}

type jwtImpl struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewJWT(cfg config.Config) JWT {
	return &jwtImpl{
		accessSecret:  []byte(cfg.JWTAccessSecret),
		refreshSecret: []byte(cfg.JWTRefreshSecret),
		accessTTL:     time.Duration(cfg.AccessTTLMin) * time.Minute,
		refreshTTL:    time.Duration(cfg.RefreshTTLHours) * time.Hour,
	}
}

func (j *jwtImpl) GenerateAccessToken(userID int64, role string) (string, time.Duration, error) {
	ttl := j.accessTTL
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString(j.accessSecret)
	return s, ttl, err
}

func (j *jwtImpl) GenerateRefreshToken(userID int64) (string, time.Duration, error) {
	ttl := j.refreshTTL
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   "refresh",
		ID:        "", // ถ้าต้องการ jti ค่อยเติม
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString(j.refreshSecret)
	return s, ttl, err
}

func (j *jwtImpl) ParseAccess(tokenStr string) (*Claims, error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return j.accessSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if c, ok := tok.Claims.(*Claims); ok && tok.Valid {
		return c, nil
	}
	return nil, errors.New("invalid token")
}

// ValidateToken is an alias for ParseAccess for middleware compatibility
func (j *jwtImpl) ValidateToken(token string) (*Claims, error) {
	return j.ParseAccess(token)
}
