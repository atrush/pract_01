package storage

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ShortURL struct {
	ID      uuid.UUID `json:"id"`
	ShortID string    `json:"shortid" validate:"required"`
	URL     string    `json:"url" validate:"required,url,max=2048"`
	UserID  uuid.UUID `json:"userid" validate:"required"`
}

type User struct {
	ID uuid.UUID `json:"id"`
}

func (s *ShortURL) Validate() error {
	validate := validator.New()

	if err := validate.Struct(s); err != nil {
		return fmt.Errorf("ошибка проверки сокращаемой ссылки: %w", err)
	}

	return nil
}
