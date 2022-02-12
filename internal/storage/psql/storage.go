package psql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/atrush/pract_01.git/internal/storage"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
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
func NewStorage(conStringDSN string, migrationsPath string) (*Storage, error) {
	if conStringDSN == "" {
		return nil, fmt.Errorf("ошибка инициализации бд:%v", "строка соединения с бд пуста")
	}

	if migrationsPath == "" {
		return nil, fmt.Errorf("ошибка инициализации бд:%v", "путь к папке миграций пуст")
	}

	db, err := sql.Open("postgres", conStringDSN)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	initMigrations(db, migrationsPath)

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
func (st *Storage) Ping() error {
	if st == nil || st.db == nil {
		return errors.New("db not initialized")
	}

	if err := st.db.Ping(); err != nil {
		return fmt.Errorf("ping for DSN (%s) failed: %w", st.conStringDSN, err)
	}

	return nil
}

// Close DB connection.
func (st Storage) Close() {
	st.db.Close()
}

// Set migrations to db
func initMigrations(db *sql.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("ошибка миграции бд:%w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return fmt.Errorf("ошибка миграции бд:%w", err)
	}

	if err = m.Up(); err != nil {
		return fmt.Errorf("ошибка миграции бд:%w", err)
	}

	return nil
}
