package schema

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/atrush/pract_01.git/internal/model"
)

type (
	//  ShortURL storage url entity.
	ShortURL struct {
		ID        uuid.UUID
		ShortID   string    `validate:"required"`
		URL       string    `validate:"required,max=2048"`
		UserID    uuid.UUID `validate:"required"`
		IsDeleted bool
	}
	//  URLList list of storage url entityes.
	URLList []ShortURL
)

//  NewURLFromCanonical creates a new ShortURL storage object from canonical model.
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

//  ToCanonical converts a storage url object to canonical model.
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

//  ToCanonical converts list of storage url object to canonical model.
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

//  Validate validates storage url object.
func (o ShortURL) Validate() error {
	if o.ID == uuid.Nil {
		return errors.New("ID не может быть nil: %v")
	}

	if o.UserID == uuid.Nil {
		return errors.New("UserID не может быть nil: %v")
	}

	if !isNotEmpty3986URL(o.ShortID) {
		return fmt.Errorf("неверное значение ShortID: %v", o.ShortID)
	}

	if !isNotEmpty3986URL(o.URL) {
		return fmt.Errorf("неверное значение URL: %v", o.URL)
	}

	return nil
}

//  isNotEmpty3986URL checks that string not empty and contains only RFC3986 symbols.
func isNotEmpty3986URL(url string) bool {
	ch := `ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789:/?#[]@!$&'()*+,;=-_.~%`

	if url == "" || len(url) > 2048 {
		return false
	}

	for _, c := range url {
		if !strings.Contains(ch, string(c)) {
			return false
		}
	}
	return true
}
