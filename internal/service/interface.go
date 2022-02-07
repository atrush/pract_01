package service

type URLShortener interface {
	GetURL(shortID string) (string, error)
	SaveURL(srcURL string, userID string) (string, error)
}

type UserManager interface {
	AddUser() (string, error)
}
