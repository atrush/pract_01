package psql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

//postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]
//const conString = "postgres://postgres:hjvfirb@localhost:5432/tst_00?sslmode=disable"

type (
	Storage struct {
		db           *bun.DB
		conStringDSN string
	}
)

// Creates a new Storage.
func NewStorage(conStringDSN string) (*Storage, error) {
	if conStringDSN == "" {
		return nil, fmt.Errorf("%s field: empty", "DSN")
	}

	conector := pgdriver.NewConnector(pgdriver.WithDSN(conStringDSN))
	sqlDB := sql.OpenDB(conector)

	st := &Storage{
		db:           bun.NewDB(sqlDB, pgdialect.New()),
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

// Close closes DB connection.
func (st Storage) Close() error {
	if st.db == nil {
		return nil
	}

	return st.db.Close()
}
