package entity

import "time"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

type Session struct {
	Token     string     `json:"token"`
	UserID    int        `json:"user_id"`
	ExpiredAt time.Time  `json:"expired_at"`
	RevokedAt *time.Time `json:"revoked_at"`
}
