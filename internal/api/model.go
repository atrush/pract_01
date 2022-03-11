package api

import (
	"errors"
	"fmt"
	"strings"

	"github.com/atrush/pract_01.git/internal/model"
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

	BatchDeleteRequest []string
)

func (s *ShortenRequest) Validate() error {
	if !IsNotEmpty3986URL(s.SrcURL) {
		return errors.New(fmt.Sprintf("неверное значение URL: %v", s.SrcURL))
	}

	return nil
}

// Make list of batch response from map[incomming-id]shotID
func NewBatchListResponseFromMap(objs map[string]string, baseURL string) []BatchResponse {
	responseArr := make([]BatchResponse, 0, len(objs))
	for k, v := range objs {
		responseArr = append(responseArr, BatchResponse{
			ID:       k,
			ShortURL: baseURL + "/" + v,
		})
	}
	return responseArr
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
