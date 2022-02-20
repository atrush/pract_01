package psql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/atrush/pract_01.git/internal/storage"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var _ storage.Storage = (*Storage)(nil)

type (
	Storage struct {
		shortURLRepo *shortURLRepository
		userRepo     *userRepository
		db           *sql.DB
		conStringDSN string
	}
)

// Create a new Storage.
func NewStorage(conStringDSN string) (*Storage, error) {
	if conStringDSN == "" {
		return nil, fmt.Errorf("ошибка инициализации бд:%v", "строка соединения с бд пуста")
	}

	db, err := sql.Open("postgres", conStringDSN)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err := initBase(db); err != nil {
		return nil, err
	}
	//initMigrations(db, migrationsPath)

	st := &Storage{
		db:           db,
		conStringDSN: conStringDSN,
	}

	st.shortURLRepo = newShortURLRepository(db)
	st.userRepo = newUserRepository(db)

	return st, nil
}

// Return URL repository
func (s *Storage) URL() storage.URLRepository {
	if s.shortURLRepo != nil {
		return s.shortURLRepo
	}
	return s.shortURLRepo
}

// Return User repository
func (s *Storage) User() storage.UserRepository {
	return s.userRepo
}

// Check DB connection.
func (s *Storage) Ping() error {
	if s == nil || s.db == nil {
		return errors.New("db not initialized")
	}

	if err := s.db.Ping(); err != nil {
		return fmt.Errorf("ping for DSN (%s) failed: %w", s.conStringDSN, err)
	}

	return nil
}

// Close DB connection.
func (s Storage) Close() {
	if s.db == nil {
		return
	}

	s.db.Close()
	s.db = nil
}

func initBase(db *sql.DB) error {
	db.Exec("DROP SCHEMA public CASCADE;CREATE SCHEMA public;")
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (" +
		"		id uuid not null,						" +
		"		primary key (id));" +
		"	CREATE TABLE IF NOT EXISTS urls(" +
		"		id uuid not null," +
		"		user_id uuid not null," +
		"		srcurl varchar(2050) not null," +
		"		shorturl varchar (16) not null," +
		"		isdeleted boolean not null," +
		"		unique (shorturl)," +
		"		unique (srcurl)," +
		"		primary key (id)," +
		"		foreign key (user_id) references users (id)" +
		"	);")
	if err != nil {
		return err
	}
	return nil

}
