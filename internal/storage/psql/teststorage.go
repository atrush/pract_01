package psql

import (
	"fmt"
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type TstStorage struct {
	Storage
}

// Test storage
func NewTestStorage(databaseURL string) (*TstStorage, error) {

	storage, err := NewStorage(databaseURL)
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
