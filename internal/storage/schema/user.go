package schema

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/atrush/pract_01.git/internal/model"
)

type (
	User struct {
		ID uuid.UUID `validate:"required"`
	}
)

// NewOrderFromCanonical creates a new ShortURL DB object from canonical model.
func NewUserFromCanonical(obj model.User) (User, error) {
	dbObj := User{
		ID: obj.ID,
	}
	if err := dbObj.Validate(); err != nil {
		return User{}, err
	}
	return dbObj, nil
}

// ToCanonical converts a DB object to canonical model.
func (u User) ToCanonical() (model.User, error) {
	obj := model.User{
		ID: u.ID,
	}

	if err := obj.Validate(); err != nil {
		return model.User{}, fmt.Errorf("status: %w", err)
	}

	return obj, nil
}

// Validate validate db obj
func (u User) Validate() error {
	validate := validator.New()

	if err := validate.Struct(u); err != nil {
		return fmt.Errorf("error validation db User : %w", err)
	}

	return nil
}
