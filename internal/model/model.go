package model

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type ShortURL struct {
	ID        uuid.UUID `json:"id"`
	ShortID   string    `json:"shortid"`
	URL       string    `json:"url"`
	UserID    uuid.UUID `json:"userid"`
	IsDeleted bool      `json:"isdeleted"`
}

type ShortURLValidator func(u ShortURL) error

func NewShortURL(srcURL string, userID uuid.UUID) ShortURL {
	return ShortURL{
		ID:        uuid.New(),
		URL:       srcURL,
		UserID:    userID,
		IsDeleted: false,
	}
}

func (u ShortURL) Validate(opts ...ShortURLValidator) error {

	if u.ID == uuid.Nil {
		return errors.New("ID не может быть nil: %v")
	}

	if u.UserID == uuid.Nil {
		return errors.New("UserID не может быть nil: %v")
	}

	if !IsNotEmpty3986URL(u.ShortID) {
		return fmt.Errorf("неверное значение ShortID: %v", u.ShortID)
	}

	if !IsNotEmpty3986URL(u.URL) {
		return fmt.Errorf("неверное значение URL: %v", u.URL)
	}

	for _, opt := range opts {
		if err := opt(u); err != nil {
			return err
		}
	}
	return nil
}

type User struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

func NewUser() User {
	return User{
		ID: uuid.New(),
	}
}

func (u User) Validate() error {
	if u.ID == uuid.Nil {
		return errors.New("ID не может быть nil: %v")
	}
	return nil
}

func IsNotEmpty3986URL(url string) bool {
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
