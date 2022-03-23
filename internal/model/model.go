package model

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

//  ShortURL represents stored url.
type ShortURL struct {
	ID        uuid.UUID `json:"id"`
	ShortID   string    `json:"shortid"`
	URL       string    `json:"url"`
	UserID    uuid.UUID `json:"userid"`
	IsDeleted bool      `json:"isdeleted"`
}

//  ShortURL rule for short url validation.
type ShortURLValidator func(u ShortURL) error

//  NewShortURL returns new ShortURL object.
//  Inits without shortID
func NewShortURL(srcURL string, userID uuid.UUID) ShortURL {
	return ShortURL{
		ID:        uuid.New(),
		URL:       srcURL,
		UserID:    userID,
		IsDeleted: false,
	}
}

//  Validate validates ShortURL object.
func (u ShortURL) Validate(opts ...ShortURLValidator) error {

	if u.ID == uuid.Nil {
		return errors.New("ID не может быть nil")
	}

	if u.UserID == uuid.Nil {
		return errors.New("UserID не может быть nil: %v")
	}

	if !isNotEmpty3986URL(u.ShortID) {
		return fmt.Errorf("неверное значение ShortID: %v", u.ShortID)
	}

	if !isNotEmpty3986URL(u.URL) {
		return fmt.Errorf("неверное значение URL: %v", u.URL)
	}

	for _, opt := range opts {
		if err := opt(u); err != nil {
			return err
		}
	}
	return nil
}

//  ShortURL represents stored user.
type User struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

// NewUser returns new NewUser object.
func NewUser() User {
	return User{
		ID: uuid.New(),
	}
}

//  Validate validates User object.
func (u User) Validate() error {
	if u.ID == uuid.Nil {
		return errors.New("ID не может быть nil: %v")
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
