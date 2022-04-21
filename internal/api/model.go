package api

import (
	"fmt"
	"strings"

	"github.com/atrush/pract_01.git/internal/model"
)

type (
	//  ShortenRequest request to save the link.
	ShortenRequest struct {
		SrcURL string `json:"url" validate:"required,url"`
	}

	//  ShortenRequest response with shorten url.
	ShortenResponse struct {
		Result string `json:"result"`
	}

	//  ShortenListResponse response list item with shorten url.
	ShortenListResponse struct {
		ShortURL string `json:"short_url"`
		SrcURL   string `json:"original_url"`
	}

	//  BatchRequest request item of list links to save, with external id.
	BatchRequest struct {
		ID  string `json:"correlation_id"`
		URL string `json:"original_url"`
	}

	//  BatchResponse response list item with external id and shorten url.
	BatchResponse struct {
		ID       string `json:"correlation_id"`
		ShortURL string `json:"short_url"`
	}

	//  BatchDeleteRequest request array of urls to delete.
	BatchDeleteRequest []string

	//  StatsResponse response stats of stored users and not deleted urls.
	StatsResponse struct {
		Urls  int `urls`
		Users int `users`
	}
)

func (s *ShortenRequest) Validate() error {
	if !isNotEmpty3986URL(s.SrcURL) {
		return fmt.Errorf("неверное значение URL: %v", s.SrcURL)
	}

	return nil
}

// NewBatchListResponseFromMap makes list of batch response from map[incoming-id]shotID.
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

//  NewShortenListResponseFromCanonical makes list of short response from arr of canonical URLs.
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
