package service

import (
	"errors"

	"github.com/atrush/pract_01.git/internal/storage"
	"github.com/google/uuid"
)

var _ UserManager = (*UserService)(nil)

type UserService struct {
	db storage.Storage
}

// New user service
func NewUserService(db storage.Storage) (*UserService, error) {
	if db == nil {
		return nil, errors.New("ошибка инициализации хранилища")
	}

	return &UserService{
		db: db,
	}, nil
}

// Check user is exist
func (u *UserService) Exist(id uuid.UUID) bool {
	if id != uuid.Nil {
		return u.db.User().Exist(id)
	}

	return false
}

// Add new user
func (u *UserService) AddUser() (*storage.User, error) {
	newUser := storage.User{
		ID: uuid.New(),
	}

	if err := u.db.User().AddUser(&newUser); err != nil {
		return nil, err
	}

	return &newUser, nil
}
