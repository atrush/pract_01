package storage

type ShortURL struct {
	ShortID string `json:"shortid"`
	URL     string `json:"url"`
	UserID  string `json:"userid"`
}
