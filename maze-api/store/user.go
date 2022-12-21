package store

import (
	"github.com/egurnov/maze-api/maze-api/model"
)

var _ model.UserStore = &UserStore{&Store{}}

type UserStore struct{ *Store }

type User struct {
	ID           int64  `json:"id,omitempty" gorm:"primary_key;auto_increment"`
	Username     string `json:"username,omitempty" gorm:"unique;not null;type:varchar(100)"`
	PasswordHash string `json:"-" gorm:"not null;type:varchar(100)"`

	Mazes []*Maze `json:"omitempty"`
}

func (s *UserStore) GetByID(id int64) (*model.User, error) {
	var user User
	err := s.db.First(&user, id).Error
	return &model.User{
		ID:           user.ID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	}, wrapError(err)
}

func (s *UserStore) GetByUsername(username string) (*model.User, error) {
	var user User
	err := s.db.First(&user, &model.User{Username: username}).Error
	return &model.User{
		ID:           user.ID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	}, wrapError(err)
}

func (s *UserStore) Create(user *model.User) (int64, error) {
	dbUser := User{
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	}
	err := s.db.Create(&dbUser).Error

	return dbUser.ID, wrapError(err)
}
