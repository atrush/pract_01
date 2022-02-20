package api

import (
	"fmt"

	"github.com/atrush/pract_01.git/internal/model"
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

// Make list of batch response from map[incomming-id]shotID
func NewBatchListResponseFromMap(objs map[string]string, baseURL string) []BatchResponse {
	var respArr []BatchResponse
	for k, v := range objs {
		respArr = append(respArr, BatchResponse{
			ID:       k,
			ShortURL: baseURL + "/" + v,
		})
	}
	return respArr
}

// Make list of short response from arr of canonical URLs
func NewShortenListResponseFromCanonical(objs []model.ShortURL, baseURL string) []ShortenListResponse {
	responseArr := make([]ShortenListResponse, 0, len(objs))
	for _, v := range objs {
		responseArr = append(responseArr, ShortenListResponse{
			ShortURL: baseURL + "/" + v.ShortID,
			SrcURL:   v.URL,
		})
	}
	return responseArr
}
