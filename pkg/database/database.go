package database

import (
	"beneburg/pkg/database/model"
	"beneburg/pkg/database/query"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

//go:generate mockgen -source=database.go -destination=./mocks/mock_database.go -package=mock_database
type Database interface {
	AutoMigrate(models ...interface{}) error
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	UpdateOrCreateUser(ctx context.Context, user *model.User) (*model.User, error)

	// CreateToken creates a new token for the given telegramID and returns token's uuid.
	CreateToken(ctx context.Context, telegramID int64) (*model.Token, error)

	GetAllUsers(ctx context.Context) ([]*model.User, error)
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*model.User, error)
	UpdateUserByID(ctx context.Context, id uint, user *model.User, updateFieldNames ...string) (*gen.ResultInfo, error)
	GetUserIDByToken(ctx context.Context, token string) (uint, error)
}

var Models = []interface{}{model.User{}, model.Token{}}

type database struct {
	db *gorm.DB

	logger *zap.Logger
}

var _ Database = database{}

func (d database) UpdateOrCreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	q := query.Use(d.db)
	u := q.User
	err := u.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "telegram_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"username"}),
	}).Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d database) GetUserIDByToken(ctx context.Context, token string) (uint, error) {
	q := query.Use(d.db)
	t := q.Token
	u := q.User
	userId, err := u.WithContext(ctx).Select(u.ID).Join(t, u.TelegramID.EqCol(t.UserTelegramId)).Where(t.UUID.Eq(token)).Where(t.ExpireAt.GtCol(t.ExpireAt.Now())).Take()
	if err != nil {
		return 0, err
	}
	return userId.ID, nil
}

func (d database) CreateToken(ctx context.Context, telegramID int64) (*model.Token, error) {
	var token *model.Token
	q := query.Use(d.db)
	t := q.Token

	err := q.Transaction(func(tx *query.Query) error {
		uid := uuid.New().String()
		err := tx.Token.WithContext(ctx).Create(&model.Token{
			UUID:           uid,
			UserTelegramId: telegramID,
			ExpireAt:       time.Now().Add(time.Hour * 24),
		})
		if err != nil {
			return err
		}
		token, err = tx.Token.WithContext(ctx).Where(t.UUID.Eq(uid)).First()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
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

// UpdateUserByID updates user by id.
// updateFieldNames is a list of fields to update, if empty all fields will be updated.
func (d database) UpdateUserByID(ctx context.Context, id uint, user *model.User, updateFieldNames ...string) (*gen.ResultInfo, error) {
	u := query.Use(d.db).User
	userDo := u.WithContext(ctx)
	if len(updateFieldNames) > 0 {
		updateFields := make([]field.Expr, 0, len(updateFieldNames))
		for _, name := range updateFieldNames {
			got, ok := u.GetFieldByName(name)
			if !ok {
				return nil, fmt.Errorf("field %s not found", name)
			}
			updateFields = append(updateFields, got)
		}
		userDo = userDo.Select(updateFields...)
	}
	update, err := userDo.Where(u.ID.Eq(id)).Updates(user)
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
