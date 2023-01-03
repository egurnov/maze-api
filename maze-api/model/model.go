package model

import (
	"context"
	"errors"
)

type User struct {
	ID int64 `json:"id,omitempty"`

	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	PasswordHash string `json:"-"`

	// Has
	Mazes []*Maze `json:"omitempty"`
}

type Maze struct {
	ID int64

	Rows     int
	Cols     int
	Entrance string
	Walls    []string

	// Foreign key
	UserID int64 `json:"-"`
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

type MazeStore interface {
	GetByID(id, userId int64) (*Maze, error)
	GetAll(userId int64) ([]*Maze, error)
	Create(*Maze) (int64, error)
	Close() error
}

type MazeService interface {
	GetByID(id, userId int64) (*Maze, error)
	GetAll(userId int64) ([]*Maze, error)
	Create(*Maze) (int64, error)
	Solve(ctx context.Context, id, userId int64, steps string) ([]string, error)
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
	ErrorNoSolution        = errors.New("no solution")
	ErrorTimelimitReached  = errors.New("time limit reached")
)
