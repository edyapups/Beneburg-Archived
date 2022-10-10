package database

import (
	"context"
	"database/sql"
	"go.uber.org/zap"
)

import _ "github.com/go-sql-driver/mysql"

type Database interface {
	Close() error
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
