package psql

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

type (
	Storage struct {
		db           *sql.DB
		conStringDSN string
	}
)

// Creates a new Storage.
func NewStorage(conStringDSN string) (*Storage, error) {
	if conStringDSN == "" {
		return nil, fmt.Errorf("%s field: empty", "DSN")
	}

	db, err := sql.Open("postgres", conStringDSN)
	if err != nil {
		return nil, err
	}

	st := &Storage{
		db:           db,
		conStringDSN: conStringDSN,
	}

	st.Ping()
	if err := st.Ping(); err != nil {
		return nil, err
	}

	return st, nil
}

// Check DB connection.
func (st *Storage) Ping() error {
	if st.db == nil {
		return errors.New("db not initialized")
	}
	if err := st.db.Ping(); err != nil {
		return fmt.Errorf("ping for DSN (%s) failed: %w", st.conStringDSN, err)
	}

	return nil
}

// Close DB connection.
func (st Storage) Close() error {
	if st.db == nil {
		return nil
	}

	return st.db.Close()
}
