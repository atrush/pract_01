package storage

import (
	"github.com/google/uuid"
)

type ShortURL struct {
	ID      uuid.UUID `json:"id"`
	ShortID string    `json:"shortid"`
	URL     string    `json:"url"`
	UserID  uuid.UUID `json:"userid"`
}

type User struct {
	ID uuid.UUID `json:"id"`
}
