package api

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type (
	ShortenRequest struct {
		SrcURL string `json:"url" validate:"required,url"`
	}

	ShortenResponse struct {
		Result string `json:"result"`
	}

	ShortenListResponse struct {
		ShortURL string `json:"short_url"`
		SrcURL   string `json:"original_url"`
	}

	BatchRequest struct {
		ID  string `json:"correlation_id"`
		URL string `json:"original_url"`
	}

	BatchResponse struct {
		ID       string `json:"correlation_id"`
		ShortURL string `json:"short_url"`
	}
)

func (s *ShortenRequest) Validate() error {
	validate := validator.New()

	if err := validate.Struct(s); err != nil {
		return fmt.Errorf("ошибка проверки сокращаемой ссылки: %w", err)
	}

	return nil
}
