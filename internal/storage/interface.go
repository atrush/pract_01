package storage

type Storage interface {
	URL() URLRepository
	User() UserRepository
	Close()
}

type URLRepository interface {
	GetURL(shortID string) (string, error)
	GetUserURLList(userID string) ([]ShortURL, error)
	SaveURL(shortID string, srcURL string, userID string) (string, error)
	IsAvailableID(shortID string) bool
}

type UserRepository interface {
	AddUser(userID string) error
	IsAvailableUserID(userID string) bool
}
