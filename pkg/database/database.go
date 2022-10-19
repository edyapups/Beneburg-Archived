package database

import (
	"beneburg/pkg/database/model"
	"beneburg/pkg/database/query"
	"context"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

//go:generate mockgen -source=database.go -destination=./mocks/mock_database.go -package=mock_database
type Database interface {
	AutoMigrate(models ...interface{}) error
	CreateUser(user *model.User) error

	// CreateToken creates a new token for the given telegramID and returns token's uuid.
	CreateToken(telegramID int64) (string, error)

	// GenerateCode generates gorm code for the given models.
	GenerateCode(models ...interface{})
	GetAllUsers() ([]*model.User, error)
	GetUserByID(id uint) (*model.User, error)
	GetUserByTelegramID(telegramID int64) (*model.User, error)
	UpdateUserByID(id uint, user *model.User) (*gen.ResultInfo, error)
	UpdateUserByTelegramID(telegramID int64, user *model.User) (*gen.ResultInfo, error)
}

type database struct {
	db *gorm.DB

	ctx    context.Context
	logger *zap.Logger
}

func (d database) CreateToken(telegramID int64) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (d database) AutoMigrate(models ...interface{}) error {
	err := d.db.AutoMigrate(models...)
	if err != nil {
		return err
	}
	return nil
}

func (d database) CreateUser(user *model.User) error {
	u := query.Use(d.db).User
	err := u.WithContext(d.ctx).Create(user)
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
func (d database) GetAllUsers() ([]*model.User, error) {
	u := query.Use(d.db).User
	all, err := u.WithContext(d.ctx).Find()
	if err != nil {
		return nil, err
	}
	return all, nil
}

func (d database) GetUserByID(id uint) (*model.User, error) {
	u := query.Use(d.db).User
	first, err := u.WithContext(d.ctx).Where(u.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}

	return first, nil
}

func (d database) GetUserByTelegramID(telegramID int64) (*model.User, error) {
	u := query.Use(d.db).User
	first, err := u.WithContext(d.ctx).Where(u.TelegramID.Eq(telegramID)).First()
	if err != nil {
		return nil, err
	}

	return first, nil
}

func (d database) UpdateUserByID(id uint, user *model.User) (*gen.ResultInfo, error) {
	u := query.Use(d.db).User
	update, err := u.WithContext(d.ctx).Update(u.ID.Eq(id), user)
	if err != nil {
		return nil, err
	}

	return &update, nil
}
func (d database) UpdateUserByTelegramID(telegramID int64, user *model.User) (*gen.ResultInfo, error) {
	u := query.Use(d.db).User
	update, err := u.WithContext(d.ctx).Update(u.TelegramID.Eq(telegramID), user)
	if err != nil {
		return nil, err
	}

	return &update, nil
}

func NewDatabase(ctx context.Context, dsn string, logger *zap.Logger) (Database, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &database{db, ctx, logger}, nil
}

func NewDatabaseWithDb(ctx context.Context, db *gorm.DB, logger *zap.Logger) Database {
	return &database{db, ctx, logger}
}
