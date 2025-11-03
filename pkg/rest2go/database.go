package rest2go

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/pressly/goose/v3"
	"github.com/ternaryss/rest2go/pkg/rest2go/settings"
)

var (
	loadOnce sync.Once
	cached   *dbProvider
	initErr  error
)

type DbCtx struct {
	Tx *sql.Tx
}

func NewDbContext(tx *sql.Tx) *DbCtx {
	return &DbCtx{Tx: tx}
}

type DbStore interface {
	Begin() (*DbCtx, error)
	Commit(context *DbCtx) error
	Rollback(context *DbCtx) error
}

type dbProvider struct {
	conf settings.Database
	db   *sql.DB
}

func NewDbProvider(conf settings.Database) (*dbProvider, error) {
	loadOnce.Do(func() {
		var db *sql.DB

		switch conf.Driver {
		case "sqlite3":
			db, initErr = initSQLiteConnection(conf)

		case "postgres":
			db, initErr = initPostgresConnection(conf)

		default:
			initErr = fmt.Errorf("unsupported database driver: %s", conf.Driver)
		}

		if initErr == nil {
			cached = &dbProvider{
				conf: conf,
				db:   db,
			}
		}
	})

	if initErr != nil {
		return nil, initErr
	}

	return cached, nil
}

func (p *dbProvider) Db() *sql.DB {
	return p.db
}

func (p *dbProvider) CloseConnection() error {
	return p.db.Close()
}

func (p *dbProvider) MigrateDatabase() error {
	migrations := fmt.Sprintf("./migrations/%s", p.conf.Driver)

	if err := goose.SetDialect(p.conf.Driver); err != nil {
		return err
	}

	if err := goose.Up(p.db, migrations); err != nil {
		return err
	}

	return nil
}

func initSQLiteConnection(conf settings.Database) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", conf.Host)

	if err != nil {
		return nil, err
	}

	if _, err := db.Exec("PRAGMA FOREIGN_KEYS=ON;"); err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func initPostgresConnection(conf settings.Database) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable&search_path=%s",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Name,
		conf.Schema,
	)
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
