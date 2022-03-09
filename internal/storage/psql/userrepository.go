package psql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/atrush/pract_01.git/internal/model"
	"github.com/atrush/pract_01.git/internal/storage"
	"github.com/google/uuid"
)

var _ storage.UserRepository = (*userRepository)(nil)

type userRepository struct {
	db *sql.DB
}

// New postgress user repository
func newUserRepository(db *sql.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

// Add user to db
func (r *userRepository) AddUser(ctx context.Context, user *model.User) error {
	if user == nil {
		return errors.New("user is nil")
	}

	result := uuid.Nil

	return r.db.QueryRowContext(
		ctx,
		"INSERT INTO users (id) VALUES ($1) RETURNING id",
		user.ID,
	).Scan(&result)
}

// Check userID exist
func (r *userRepository) Exist(userID uuid.UUID) (bool, error) {
	count := 0
	err := r.db.QueryRow(
		"SELECT  COUNT(*) as count FROM users WHERE id = $1", userID).Scan(&count)

	if err != nil {
		return false, err
	}
	return count > 0, nil
}
