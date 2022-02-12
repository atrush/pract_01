package psql

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	st "github.com/atrush/pract_01.git/internal/storage"
)

var _ st.URLRepository = (*shortURLRepository)(nil)

type shortURLRepository struct {
	db *sql.DB
}

// New postgress URL repository
func newShortURLRepository(db *sql.DB) *shortURLRepository {
	return &shortURLRepository{
		db: db,
	}
}

// Save URL to db
func (r *shortURLRepository) SaveURL(sht *st.ShortURL) error {
	if err := sht.Validate(); err != nil {
		return fmt.Errorf("ошибка хранилица:%w", err)
	}

	return r.db.QueryRow(
		"INSERT INTO urls (id, user_id, srcurl, shorturl) VALUES ($1, $2, $3, $4) RETURNING id",
		sht.ID,
		sht.UserID,
		sht.URL,
		sht.ShortID,
	).Scan(&sht.ID)
}

// Get source URL by shortID from db
func (r *shortURLRepository) GetURL(shortID string) (string, error) {
	res := ""
	err := r.db.QueryRow(
		"select srcurl from urls where shorturl = $1", shortID,
	).Scan(&res)

	if err != nil {
		return "", fmt.Errorf("ошибка хранилица:%w", err)
	}

	return res, nil
}

// Get users urls by user id
func (r *shortURLRepository) GetUserURLList(userID uuid.UUID, limit int) ([]st.ShortURL, error) {
	userURLs := make([]st.ShortURL, 0, limit)

	rows, err := r.db.Query(
		"SELECT id, user_id, srcurl, shorturl from urls WHERE user_id = $1 LIMIT $2", userID, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var s st.ShortURL
		err = rows.Scan(&s.ID, &s.UserID, &s.URL, &s.ShortID)
		if err != nil {
			return nil, err
		}

		userURLs = append(userURLs, s)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return userURLs, nil
}

// check shortID exist in db
func (r *shortURLRepository) Exist(shortID string) (bool, error) {
	count := 0
	err := r.db.QueryRow(
		"SELECT  COUNT(*) as count FROM urls WHERE shorturl = $1", shortID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
