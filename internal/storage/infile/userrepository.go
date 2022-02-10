package infile

import (
	"errors"
	"sync"

	"github.com/atrush/pract_01.git/internal/storage"
	"github.com/google/uuid"
)

var _ storage.UserRepository = (*userRepository)(nil)

type userRepository struct {
	cache cache
	sync.RWMutex
}

// Init new repository
func newUserRepository(c *cache) (*userRepository, error) {
	if c == nil {
		return nil, errors.New("cant init repository cache not init")
	}

	return &userRepository{
		cache: *c,
	}, nil
}

// Check userID exist
func (r *userRepository) Exist(userID uuid.UUID) bool {
	r.RLock()
	_, ok := r.cache.userCache[userID]
	defer r.RUnlock()

	return ok
}

// Add User
func (r *userRepository) AddUser(user *storage.User) error {
	r.Lock()
	r.cache.userCache[user.ID] = user.ID
	defer r.Unlock()

	return nil
}
