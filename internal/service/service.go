package service

import (
	"errors"
	"fmt"

	st "github.com/atrush/pract_01.git/internal/storage"
)

var _ Servicer = (*Service)(nil)

type Service struct {
	shotURLService *ShortURLService
	userService    *UserService
}

// New Service
func NewService(db st.Storage) (*Service, error) {
	if db == nil {
		return nil, errors.New("ошибка инициализации сервиса: хранилище nil")
	}

	shortURLService, err := newShortURLService(db)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации сервиса:%w", err)
	}

	userService, err := newUserService(db)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации сервиса:%w", err)
	}
	return &Service{
		shotURLService: shortURLService,
		userService:    userService,
	}, nil
}

// Return URL service
func (s *Service) URL() URLShortener {
	return s.shotURLService
}

// Return User service
func (s *Service) User() UserManager {
	return s.userService
}
