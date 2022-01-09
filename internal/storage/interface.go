package storage

type URLStorage struct {
	URLStorage URLStorer
}
type URLStorer interface {
	GetURL(shortURL string) (string, error)
	SaveURL(srcURL string) (string, error)
}

func NewURLStorage(db URLStorer) *URLStorage {
	return &URLStorage{
		URLStorage: db,
	}
}
