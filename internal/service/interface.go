package service

type URLShortener interface {
	GetURL(shortID string) (string, error)
	SaveURL(srcURL string) (string, error)
}
