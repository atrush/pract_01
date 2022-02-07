package storage

type URLStorer interface {
	GetURL(shortID string) (string, error)
	SaveURL(shortID string, srcURL string, userID string) (string, error)
	IsAvailableID(shortID string) bool
}

type UserStorer interface {
	AddUser(userID string) error
	IsAvailableUserID(userID string) bool
}
