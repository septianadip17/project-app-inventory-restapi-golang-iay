package service

import (
	"context"
	"errors"
	"project-app-inventory-restapi-golang-iay/internal/entity"
	"project-app-inventory-restapi-golang-iay/internal/repository"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	token := uuid.NewString()
	session := entity.Session{
		Token:     token,
		UserID:    user.ID,
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return "", err
	}

	return token, nil
}
