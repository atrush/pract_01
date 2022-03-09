package service

import (
	"context"
	"errors"

	"github.com/atrush/pract_01.git/internal/model"
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
func (u *UserService) Exist(ctx context.Context, id uuid.UUID) (bool, error) {
	if id == uuid.Nil {
		return false, errors.New("ошибка проверки существования user: uuid nil")
	}

	return u.db.User().Exist(id)
}

// Add new user
func (u *UserService) AddUser(ctx context.Context) (*model.User, error) {
	newUser := model.User{
		ID: uuid.New(),
	}

	if err := u.db.User().AddUser(ctx, &newUser); err != nil {
		return nil, err
	}

	return &newUser, nil
}
