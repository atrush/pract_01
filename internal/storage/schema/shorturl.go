package schema

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/atrush/pract_01.git/internal/model"
)

type (
	ShortURL struct {
		ID        uuid.UUID
		ShortID   string    `validate:"required"`
		URL       string    `validate:"required,max=2048"`
		UserID    uuid.UUID `validate:"required"`
		IsDeleted bool
	}
	URLList []ShortURL
)

// NewOrderFromCanonical creates a new ShortURL DB object from canonical model.
func NewURLFromCanonical(obj model.ShortURL) (ShortURL, error) {
	dbObj := ShortURL{
		ID:        obj.ID,
		ShortID:   obj.ShortID,
		URL:       obj.URL,
		UserID:    obj.UserID,
		IsDeleted: obj.IsDeleted,
	}
	if err := dbObj.Validate(); err != nil {
		return ShortURL{}, err
	}
	return dbObj, nil
}

// ToCanonical converts a DB object to canonical model.
func (o ShortURL) ToCanonical() (model.ShortURL, error) {
	obj := model.ShortURL{
		ID:        o.ID,
		ShortID:   o.ShortID,
		URL:       o.URL,
		UserID:    o.UserID,
		IsDeleted: o.IsDeleted,
	}

	if err := obj.Validate(); err != nil {
		return model.ShortURL{}, fmt.Errorf("status: %w", err)
	}

	return obj, nil
}

// ToCanonical converts a DB object to canonical model.
func (o URLList) ToCanonical() ([]model.ShortURL, error) {
	objs := make([]model.ShortURL, 0, len(o))
	for dbObjIdx, dbObj := range o {
		obj, err := dbObj.ToCanonical()
		if err != nil {
			return nil, fmt.Errorf("dbObj [%d]: %w", dbObjIdx, err)
		}
		objs = append(objs, obj)
	}

	return objs, nil
}

// Validate validate db obj
func (o ShortURL) Validate() error {
	validate := validator.New()

	if err := validate.Struct(o); err != nil {
		return fmt.Errorf("error validation db ShortURL : %w", err)
	}

	return nil
}
