package service

import (
	"context"
	"errors"
	"project-app-inventory-restapi-golang-iay/internal/entity"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	// 1. ambil username
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// 2. Cek Password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// 3. buat session uuid
	token := uuid.NewString()
	session := entity.Session{
		Token:     token,
		UserID:    user.ID,
		ExpiredAt: time.Now().Add(24 * time.Hour), // Expire 1 hari
	}

	// 4. Simpan ke DB
	err = s.sessionRepo.CreateSession(ctx, session)
	return token, err
}
