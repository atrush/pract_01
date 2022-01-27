package api

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ShortenRequest struct {
	SrcURL string `json:"url" validate:"required,url"`
}

func (s *ShortenRequest) Validate() error {
	validate := validator.New()

	if err := validate.Struct(s); err != nil {
		return fmt.Errorf("ошибка проверки сокращаемой ссылки: %w", err)
	}

	return nil
}

type ShortenResponse struct {
	Result string `json:"result"`
}
