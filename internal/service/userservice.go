package service

import (
	"errors"

	"github.com/atrush/pract_01.git/internal/storage"
)

var _ UserManager = (*UserService)(nil)

type UserService struct {
	db storage.Storage
}

func (u *UserService) UserExist(userID string) bool {

	return !u.db.User().IsAvailableUserID(userID)
}

func (u *UserService) AddUser() (string, error) {
	newUserID, err := u.GenUserID()
	if err != nil {
		return "", err
	}

	if err := u.db.User().AddUser(newUserID); err != nil {
		return "", err
	}

	return newUserID, nil
}

func NewUserService(db storage.Storage) (*UserService, error) {
	if db == nil {
		return nil, errors.New("ошибка инициализации хранилища")
	}

	return &UserService{
		db: db,
	}, nil
}

func (u *UserService) GenUserID() (string, error) {
	dst, err := u.iterUserIDGenerator(0)
	if err != nil {
		return "", err
	}

	return dst, nil
}

func (u *UserService) iterUserIDGenerator(iterationCount int) (string, error) {
	userID, err := GenUserID()
	if err != nil {
		return "", err
	}
	if !u.db.User().IsAvailableUserID(userID) {
		iterationCount++

		userID, err := u.iterUserIDGenerator(iterationCount)
		if err != nil || iterationCount > 10 {
			return "", errors.New("ошибка генерации userID")
		}

		return userID, nil
	}

	return userID, nil
}
