package psql

import (
	"fmt"
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const (
	databaseURL    = "postgres://postgres:hjvfirb@localhost:5432/tst?sslmode=disable"
	migrationsPath = "file://migrations"
)

type TstStorage struct {
	Storage
}

// Test storage
func NewTestStorage() (*TstStorage, error) {

	storage, err := NewStorage(databaseURL, migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициации тестовой бд:%w", err)
	}

	return &TstStorage{Storage: *storage}, nil
}

// Clear db
func (t *TstStorage) DropAll() {
	_, err := t.db.Exec("DROP SCHEMA public CASCADE;CREATE SCHEMA public;")

	if err != nil {
		log.Fatal("ошибка очистки бд:" + err.Error())
	}
}
