package database

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

type Database interface {
	AutoMigrate(models ...interface{}) error
	GenerateCode(models ...interface{})
}

type database struct {
	db *gorm.DB

	ctx    context.Context
	logger *zap.Logger
}

func (d database) AutoMigrate(models ...interface{}) error {
	err := d.db.AutoMigrate(models...)
	if err != nil {
		return err
	}
	return nil
}

func (d database) GenerateCode(models ...interface{}) {
	g := gen.NewGenerator(gen.Config{
		OutPath:           "pkg/database/query",
		FieldNullable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})
	g.ApplyBasic(models...)
	g.Execute()
}

func NewDatabase(ctx context.Context, dsn string, logger *zap.Logger) (Database, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &database{db, ctx, logger}, nil
}
