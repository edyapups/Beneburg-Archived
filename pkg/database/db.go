package database

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

import _ "github.com/go-sql-driver/mysql"

type Database interface {
	Close() error
	PingContext(ctx context.Context) error
	MakeMigrations(forceMigration int) error
}

type database struct {
	db     *sql.DB
	ctx    context.Context
	logger *zap.Logger
}

func NewDatabase(ctx context.Context, logger *zap.Logger, dataSourceName string) (Database, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &database{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}, nil
}

func (d *database) Close() error {
	return d.db.Close()
}

func (d *database) PingContext(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

func (d *database) MakeMigrations(forceMigration int) error {
	driver, err := mysql.WithInstance(d.db, &mysql.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && forceMigration != 0 {
		err = m.Force(forceMigration)
		return err
	}
	return nil
}
