package model

import "errors"

type User struct {
	ID           int64  `json:"id,omitempty" binding:"isdefault"`
	Username     string `json:"username,omitempty" binding:"omitempty,username"`
	Password     string `json:"password,omitempty" binding:"omitempty,min=6"`
	PasswordHash string `json:"-"`
}

type CustomClaims struct {
	UserID int64 `json:"user_id"`
}

type UserStore interface {
	GetByID(id int64) (*User, error)
	GetByUsername(username string) (*User, error)
	Create(*User) (int64, error)
	Close() error
}

type UserService interface {
	GetByID(id int64) (*User, error)
	GetByUsername(username string) (*User, error)
	Create(*User) (int64, error)
}

type JWTService interface {
	GenerateToken(id int64) (string, error)
	ValidateToken(token string) (*CustomClaims, error)
}

var (
	ErrUsernameAlreadyUsed = errors.New("username already used")
	ErrInvalidInput        = errors.New("bad input")
	ErrNotFound            = errors.New("not found")
	ErrNotAllowed          = errors.New("not allowed")
	ErrorUnauthorized      = errors.New("unauthorized")
)
