package store

import (
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //nolint:golint
	"github.com/pkg/errors"

	"github.com/egurnov/maze-api/maze-api/model"
)

type Store struct {
	db *gorm.DB
}

func NewMySQLStore(filename string) (*Store, error) {
	db, err := gorm.Open("mysql", filename)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open mysql database")
	}

	db.AutoMigrate(
		&User{},
	)
	db.LogMode(false)

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Wipe() error {
	err := s.db.DropTable(&User{}).Error
	if err != nil {
		return err
	}
	err = s.db.AutoMigrate(&User{}).Error
	if err != nil {
		return err
	}

	return nil
}

func wrapError(err error) error {
	if gorm.IsRecordNotFoundError(err) {
		return model.ErrNotFound
	}
	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed: users.username") {
		return model.ErrUsernameAlreadyUsed
	}
	return err
}
