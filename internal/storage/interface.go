package storage

type URLStorer interface {
	GetURL(shortID string) (string, error)
	SaveURL(shortID string, srcURL string) (string, error)
	IsAvailableID(shortID string) bool
}
