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

type userRepository struct {
	cache *cache
}

// Init new repository
func newUserRepository(c *cache) (*userRepository, error) {
	if c == nil {
		return nil, errors.New("cant init repository cache not init")
	}

	return &userRepository{
		cache: c,
	}, nil
}

// Check userID exist
func (r *userRepository) Exist(userID uuid.UUID) (bool, error) {
	r.cache.RLock()
	_, ok := r.cache.userCache[userID]
	defer r.cache.RUnlock()

	return ok, nil
}

// Add User
func (r *userRepository) AddUser(_ context.Context, user *model.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	dbObj, err := schema.NewUserFromCanonical(*user)
	if err != nil {
		return err
	}
	r.cache.Lock()
	r.cache.userCache[user.ID] = dbObj.ID
	defer r.cache.Unlock()

	return nil
}
