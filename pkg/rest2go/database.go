package rest2go

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/pressly/goose/v3"
	"github.com/ternaryss/rest2go/pkg/rest2go/settings"
)

const (
	DbSQLite   string = "sqlite3"
	DbPostgres string = "postgres"
)

var (
	loadOnce sync.Once
	cached   *Database
	dbErr    error
)

type Database struct {
	conf settings.Database
	Pool *sql.DB
}

func NewDatabase(conf settings.Database) (*Database, error) {
	loadOnce.Do(func() {
		var pool *sql.DB

		switch conf.Driver {
		case DbSQLite:
			pool, dbErr = initSQLite(conf)

		case DbPostgres:
			pool, dbErr = initPostgres(conf)

		default:
			dbErr = fmt.Errorf("unsupported database driver: %s", conf.Driver)
			return
		}

		if dbErr == nil {
			cached = &Database{conf: conf, Pool: pool}
		}
	})

	return cached, dbErr
}

func (d *Database) Close() error {
	return d.Pool.Close()
}

func (d *Database) Migrate() error {
	migrations := fmt.Sprintf("./migrations/%s", d.conf.Driver)

	if err := goose.SetDialect(d.conf.Driver); err != nil {
		return err
	}

	if err := goose.Up(d.Pool, migrations); err != nil {
		return err
	}

	return nil
}

func (d *Database) Begin() (*sql.Tx, error) {
	return d.Pool.Begin()
}

func (d *Database) Commit(tx *sql.Tx) error {
	if tx == nil {
		return errors.New("nil transaction")
	}

	return tx.Commit()
}

func (d *Database) Rollback(tx *sql.Tx) error {
	if tx == nil {
		return errors.New("nil transaction")
	}

	return tx.Rollback()
}

func (d *Database) Exec(tx *sql.Tx, query string, args ...any) (sql.Result, error) {
	if tx != nil {
		return tx.Exec(query, args...)
	}

	return d.Pool.Exec(query, args...)
}

func (d *Database) Query(tx *sql.Tx, query string, args ...any) (*sql.Rows, error) {
	if tx != nil {
		return tx.Query(query, args...)
	}

	return d.Pool.Query(query, args...)
}

func (d *Database) QueryRow(tx *sql.Tx, query string, args ...any) *sql.Row {
	if tx != nil {
		return tx.QueryRow(query, args...)
	}

	return d.Pool.QueryRow(query, args...)
}

func initSQLite(conf settings.Database) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s?_foreign_keys=on", conf.Host)
	db, err := sql.Open(DbSQLite, dsn)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func initPostgres(conf settings.Database) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable&search_path=%s",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Name,
		conf.Schema,
	)
	db, err := sql.Open(DbPostgres, dsn)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
