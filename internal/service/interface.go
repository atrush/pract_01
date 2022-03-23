package service

import (
	"context"
	"github.com/atrush/pract_01.git/internal/model"
	"github.com/google/uuid"
)

//  URLShortener is the interface that wraps methods for process urls.
type URLShortener interface {
	//  GetURL returns canonical ShortURL by shortID.
	GetURL(ctx context.Context, shortID string) (model.ShortURL, error)

	//  GetUserURLList returns array of canonical ShortURL by userID. Or empty array if not founded.
	GetUserURLList(ctx context.Context, userID uuid.UUID) ([]model.ShortURL, error)

	//  SaveURL saves incoming URL and return shortID.
	SaveURL(ctx context.Context, srcURL string, userID uuid.UUID) (string, error)

	//  SaveURLList saves list of urls for user.
	//  Accept map of [ext_id]url, save and replace url by shortID.
	SaveURLList(srcArr map[string]string, userID uuid.UUID) error

	//  DeleteURLList marks list of short urls as deleted.
	DeleteURLList(userID uuid.UUID, shortIDList ...string) error

	//  Ping checks db connection.
	Ping(ctx context.Context) error
}

// UserManager is the interface that wraps methods for process users.
type UserManager interface {
	//  AddUser creates new user, save to storage and return instance.
	AddUser(ctx context.Context) (model.User, error)

	//  Exist checks user is exist, by user id.
	Exist(ctx context.Context, id uuid.UUID) (bool, error)
}
