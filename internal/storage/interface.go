package storage

type UrlStorage struct {
	UrlStorage UrlStorer
}
type UrlStorer interface {
	GetUrl(shortUrl string) (string, error)
	SaveUrl(srcUrl string) (string, error)
}

func NewUrlStorage(db UrlStorer) *UrlStorage {
	return &UrlStorage{
		UrlStorage: db,
	}
}
