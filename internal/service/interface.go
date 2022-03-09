package service

import (
	"context"
	"github.com/atrush/pract_01.git/internal/model"
	"github.com/google/uuid"
)

type URLShortener interface {
	GetURL(ctx context.Context, shortID string) (model.ShortURL, error)
	GetUserURLList(ctx context.Context, userID uuid.UUID) ([]model.ShortURL, error)
	SaveURL(ctx context.Context, srcURL string, userID uuid.UUID) (string, error)
	SaveURLList(srcArr map[string]string, userID uuid.UUID) error
	DeleteURLList(userID uuid.UUID, shortIDList ...string) error
	Ping(ctx context.Context) error
}

type UserManager interface {
	AddUser(ctx context.Context) (model.User, error)
	Exist(ctx context.Context, id uuid.UUID) (bool, error)
}
