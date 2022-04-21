package service

import (
	"context"
	"errors"

	"github.com/atrush/pract_01.git/internal/model"
	"github.com/atrush/pract_01.git/internal/storage"
	"github.com/google/uuid"
)

var _ UserManager = (*UserService)(nil)

//  UserService implements UserManager interface, provides operations with users.
type UserService struct {
	db storage.Storage
}

//  NewUserService inits and returns new user service.
func NewUserService(db storage.Storage) (*UserService, error) {
	if db == nil {
		return nil, errors.New("ошибка инициализации хранилища")
	}

	return &UserService{
		db: db,
	}, nil
}

//  Exist checks user is exist, by user id.
func (u *UserService) Exist(ctx context.Context, id uuid.UUID) (bool, error) {
	if id == uuid.Nil {
		return false, errors.New("ошибка проверки существования user: uuid nil")
	}

	return u.db.User().Exist(id)
}

//  GetCount returns count of stored users.
func (u *UserService) GetCount() (int, error) {
	return u.db.User().GetCount()
}

//  AddUser creates new user, save to storage and return instance.
func (u *UserService) AddUser(ctx context.Context) (model.User, error) {
	newUser := model.NewUser()

	newUser, err := u.db.User().AddUser(ctx, newUser)
	if err != nil {
		return model.User{}, err
	}
	return newUser, nil
}
