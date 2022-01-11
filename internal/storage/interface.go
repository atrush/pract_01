package storage

type URLStorer interface {
	GetURL(shortURL string) (string, error)
	SaveURL(srcURL string) (string, error)
}
