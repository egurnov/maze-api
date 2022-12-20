package service

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/egurnov/maze-api/maze-api/model"
)

type UserService struct {
	Store model.UserStore
}

var _ model.UserService = &UserService{}

func (s *UserService) GetByID(id int64) (*model.User, error) {
	return s.Store.GetByID(id)
}

func (s *UserService) GetByUsername(email string) (*model.User, error) {
	return s.Store.GetByUsername(email)
}

func (s *UserService) Create(user *model.User) (int64, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	user.PasswordHash = string(passHash)

	if len(user.Username) == 0 || len(user.Password) == 0 {
		return 0, model.ErrInvalidInput
	}

	return s.Store.Create(user)
}
