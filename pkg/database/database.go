package database

import (
	"beneburg/pkg/database/model"
	"beneburg/pkg/database/query"
	"context"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"time"
)

//go:generate mockgen -source=database.go -destination=./mocks/mock_database.go -package=mock_database
type Database interface {
	AutoMigrate(models ...interface{}) error
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)

	// CreateToken creates a new token for the given telegramID and returns token's uuid.
	CreateToken(ctx context.Context, telegramID int64) (*model.Token, error)
	GetAllUsers(ctx context.Context) ([]*model.User, error)
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*model.User, error)
	UpdateUserByID(ctx context.Context, id uint, user *model.User) (*gen.ResultInfo, error)
	UpdateUserByTelegramID(ctx context.Context, telegramID int64, user *model.User) (*gen.ResultInfo, error)
}

var Models = []interface{}{model.User{}, model.Token{}}

type database struct {
	db *gorm.DB

	logger *zap.Logger
}

func (d database) CreateToken(ctx context.Context, telegramID int64) (*model.Token, error) {
	t := query.Use(d.db).Token
	token := model.Token{
		UserTelegramId: telegramID,
		ExpireAt:       time.Now().Add(time.Hour * 24),
	}
	err := t.WithContext(ctx).Create(&token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (d database) AutoMigrate(models ...interface{}) error {
	err := d.db.AutoMigrate(models...)
	if err != nil {
		return err
	}
	return nil
}

func (d database) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	u := query.Use(d.db).User
	err := u.WithContext(ctx).Create(user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (d database) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	u := query.Use(d.db).User
	all, err := u.WithContext(ctx).Find()
	if err != nil {
		return nil, err
	}
	return all, nil
}

func (d database) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	u := query.Use(d.db).User
	first, err := u.WithContext(ctx).Where(u.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}

	return first, nil
}

func (d database) GetUserByTelegramID(ctx context.Context, telegramID int64) (*model.User, error) {
	u := query.Use(d.db).User
	first, err := u.WithContext(ctx).Where(u.TelegramID.Eq(telegramID)).First()
	if err != nil {
		return nil, err
	}

	return first, nil
}

func (d database) UpdateUserByID(ctx context.Context, id uint, user *model.User) (*gen.ResultInfo, error) {
	u := query.Use(d.db).User
	update, err := u.WithContext(ctx).Update(u.ID.Eq(id), user)
	if err != nil {
		return nil, err
	}

	return &update, nil
}
func (d database) UpdateUserByTelegramID(ctx context.Context, telegramID int64, user *model.User) (*gen.ResultInfo, error) {
	u := query.Use(d.db).User
	update, err := u.WithContext(ctx).Update(u.TelegramID.Eq(telegramID), user)
	if err != nil {
		return nil, err
	}

	return &update, nil
}

func NewDatabase(dsn string, logger *zap.Logger) (Database, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &database{db, logger}, nil
}

func NewDatabaseWithDb(db *gorm.DB, logger *zap.Logger) Database {
	return &database{db, logger}
}
