package infile

import (
	"context"
	"errors"

	"github.com/atrush/pract_01.git/internal/model"
	"github.com/atrush/pract_01.git/internal/storage"
	"github.com/atrush/pract_01.git/internal/storage/schema"
	"github.com/google/uuid"
)

var _ storage.UserRepository = (*userRepository)(nil)

//  userRepository implements UserRepository interface, provides actions with user records in inmemory storage.
type userRepository struct {
	cache *cache
}

//  newUserRepository inits new user repository.
func newUserRepository(c *cache) (*userRepository, error) {
	if c == nil {
		return nil, errors.New("cant init repository cache not init")
	}

	return &userRepository{
		cache: c,
	}, nil
}

//  AddUser saves user to inmemory storage.
func (r *userRepository) AddUser(_ context.Context, user model.User) (model.User, error) {
	dbObj, err := schema.NewUserFromCanonical(user)
	if err != nil {
		return model.User{}, err
	}
	r.cache.Lock()
	r.cache.userCache[user.ID] = dbObj.ID
	defer r.cache.Unlock()

	return user, nil
}

//  Exist checks that user is exist in storage.
func (r *userRepository) Exist(userID uuid.UUID) (bool, error) {
	r.cache.RLock()
	_, ok := r.cache.userCache[userID]
	defer r.cache.RUnlock()

	return ok, nil
}
