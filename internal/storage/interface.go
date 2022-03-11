package storage

import (
	"context"
	"github.com/atrush/pract_01.git/internal/model"
	"github.com/google/uuid"
)

//  Storage is the interface that wraps methods for working with the database.
type Storage interface {
	//  URL returns repository for working with urls.
	URL() URLRepository

	//  User returns repository for working with users.
	User() UserRepository

	//  Close closes storage connection.
	Close()

	//  Ping checks connection to storage.
	Ping() error
}

//  URLRepository is the interface that wraps methods for working with url records in database.
type URLRepository interface {
	//  GetURL selects url record from database by shortID and returns as canonical ShortURL.
	GetURL(ctx context.Context, shortID string) (model.ShortURL, error)

	//  GetUserURLList selects list of url records by user id and returns as array of canonical ShortURL.
	GetUserURLList(ctx context.Context, userID uuid.UUID, limit int) ([]model.ShortURL, error)

	//  SaveURL saves canonical ShortURL to database and returns saved instance.
	SaveURL(ctx context.Context, shURL model.ShortURL) (model.ShortURL, error)

	//  Exist checks than record with shot id is exist.
	Exist(shortID string) (bool, error)

	//  SaveURLBuff writes ShortURL elements to buffer. If buffer full runs SaveURLBuffFlush.
	SaveURLBuff(shURL *model.ShortURL) error

	//  SaveURLBuffFlush writes ShortURL elements from buffer to database, updates objects, cleans buffer.
	SaveURLBuffFlush() error

	//  DeleteURLBatch async updates list of urls as deleted.
	DeleteURLBatch(userID uuid.UUID, shortIDList ...string) error
}

//  UserRepository is the interface that wraps methods for working with url records in database.
type UserRepository interface {
	//  AddUser saves canonical User to database and returns saved instance.
	AddUser(ctx context.Context, user model.User) (model.User, error)

	//  Exist checks than record with id is exist in storage.
	Exist(userID uuid.UUID) (bool, error)
}
