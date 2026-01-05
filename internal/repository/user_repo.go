package repository

import (
	"context"
	"project-app-inventory-restapi-golang-iay/internal/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	query := "SELECT id, username, password, role FROM users WHERE username = $1"
	err := r.db.QueryRow(ctx, query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	return &user, err
}

func (r *UserRepository) CreateSession(ctx context.Context, session entity.Session) error {
	query := "INSERT INTO sessions (token, user_id, expired_at) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(ctx, query, session.Token, session.UserID, session.ExpiredAt)
	return err
}

func (r *UserRepository) GetSession(ctx context.Context, token string) (*entity.Session, error) {
	var s entity.Session
	query := "SELECT token, user_id, expired_at, revoked_at FROM sessions WHERE token = $1"
	err := r.db.QueryRow(ctx, query, token).Scan(&s.Token, &s.UserID, &s.ExpiredAt, &s.RevokedAt)
	return &s, err
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	var user entity.User
	query := "SELECT id, username, role FROM users WHERE id = $1"
	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.Role)
	return &user, err
}
