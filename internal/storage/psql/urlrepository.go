package psql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"

	"github.com/atrush/pract_01.git/internal/model"
	"github.com/atrush/pract_01.git/internal/shterrors"
	st "github.com/atrush/pract_01.git/internal/storage"
	"github.com/atrush/pract_01.git/internal/storage/schema"
)

var _ st.URLRepository = (*shortURLRepository)(nil)

type shortURLRepository struct {
	db        *sql.DB
	urlbuffer []model.ShortURL
}

// New postgress URL repository
func newShortURLRepository(db *sql.DB) *shortURLRepository {
	return &shortURLRepository{
		db:        db,
		urlbuffer: make([]model.ShortURL, 0, 100),
	}
}

// Save ShortURL using buffer
func (r *shortURLRepository) SaveURLBuff(sht *model.ShortURL) error {
	r.urlbuffer = append(r.urlbuffer, *sht)

	if cap(r.urlbuffer) == len(r.urlbuffer) {
		err := r.SaveURLBuffFlush()
		if err != nil {
			return fmt.Errorf("ошибка хранилица:%w", err)
		}
	}
	return nil
}

// Save ShorURLs stored in bufferr to db
func (r *shortURLRepository) SaveURLBuffFlush() error {
	if r.db == nil {
		return errors.New("ошибка транзакции сохранения: база не инициирована")
	}
	tx, err := r.db.Begin()
	if err != nil {
		return err

	}
	stmt, err := tx.Prepare("INSERT INTO urls(id, user_id, srcurl, shorturl) VALUES($1, $2, $3, $4)RETURNING id")
	if err != nil {
		return err
	}

	for _, sht := range r.urlbuffer {
		dbObj, err := schema.NewURLFromCanonical(sht)
		if err == nil {
			if stmt.QueryRow(dbObj.ID, dbObj.UserID, dbObj.URL, dbObj.ShortID).Scan(&dbObj.ID); err == nil {
				sht.ID = dbObj.ID
				continue
			}
		}

		if err = tx.Rollback(); err != nil {
			return fmt.Errorf("ошибка транзакции сохранения: транзакцию не удалось отменить:%w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ошибка транзакции сохранения:%w", err)
	}

	r.urlbuffer = r.urlbuffer[:0]
	return nil
}

// Save URL to db
func (r *shortURLRepository) SaveURL(sht *model.ShortURL) error {
	dbObj, err := schema.NewURLFromCanonical(*sht)
	if err != nil {
		return fmt.Errorf("ошибка хранилица:%w", err)
	}
	row := r.db.QueryRow(
		"INSERT INTO urls (id, user_id, srcurl, shorturl) VALUES ($1, $2, $3, $4) RETURNING id ",
		dbObj.ID,
		dbObj.UserID,
		dbObj.URL,
		dbObj.ShortID,
	)

	if row.Err() != nil {
		// check duplicate srcurl
		pqErr, ok := row.Err().(*pq.Error)
		if ok && pqErr.Code == pgerrcode.UniqueViolation && pqErr.Constraint == "urls_srcurl_key" {
			shortID, err := r.GetShortURLBySrcURL(sht.URL)
			if err != nil {
				return fmt.Errorf("ошибка добавления записи в БД, ссылка %v уже существует: ошибка получения существующей короткой ссыки: %w",
					sht.URL, err)
			}
			return &shterrors.ErrorConflictSaveURL{
				Err:           row.Err(),
				ExistShortURL: shortID,
			}
		}
	}

	if err := row.Scan(&dbObj.ID); err != nil {
		return err
	}

	sht.ID = dbObj.ID
	return nil
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

// Get source URL by shortID from db
func (r *shortURLRepository) GetShortURLBySrcURL(url string) (string, error) {
	res := ""
	err := r.db.QueryRow(
		"select shorturl from urls where srcurl = $1", url,
	).Scan(&res)

	if err != nil {
		return "", fmt.Errorf("ошибка хранилица:%w", err)
	}

	return res, nil
}

// Get users urls by user id
func (r *shortURLRepository) GetUserURLList(userID uuid.UUID, limit int) ([]model.ShortURL, error) {
	var userURLs schema.URLList
	userURLs = make([]schema.ShortURL, 0, limit)

	rows, err := r.db.Query(
		"SELECT id, user_id, srcurl, shorturl from urls WHERE user_id = $1 LIMIT $2", userID, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var s schema.ShortURL
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
	return userURLs.ToCanonical()
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
