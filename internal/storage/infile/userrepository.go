package infile

import (
	"errors"
	"sync"

	"github.com/atrush/pract_01.git/internal/storage"
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
func (r *userRepository) IsAvailableUserID(userID string) bool {
	r.RLock()
	_, ok := r.cache.userCache[userID]
	defer r.RUnlock()

	return !ok
}

// Add User
func (r *userRepository) AddUser(userID string) error {
	if userID == "" {
		return errors.New("нельзя использовать пустой id")
	}
	if !r.IsAvailableUserID(userID) {
		return errors.New("id уже существует")
	}

	r.Lock()
	r.cache.userCache[userID] = userID
	defer r.Unlock()

	return nil
}
