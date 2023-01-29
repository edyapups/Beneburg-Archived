package database

import (
	"beneburg/pkg/database/model"
	"beneburg/pkg/database/query"
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
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
	CreateOrProlongToken(ctx context.Context, telegramID int64) (*model.Token, error)
	ReissueToken(ctx context.Context, telegramID int64) (*model.Token, error)
	GetUserByToken(ctx context.Context, token string) (*model.User, error)

	GetAllUsers(ctx context.Context) ([]*model.User, error)
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*model.User, error)
	UpdateUserByID(ctx context.Context, id uint, user *model.User) (*model.User, error)

	CreateForm(ctx context.Context, form *model.Form) (*model.Form, error)
	AcceptForm(ctx context.Context, id uint) (*gen.ResultInfo, error)
	RejectForm(ctx context.Context, id uint) (*gen.ResultInfo, error)
	GetActualForm(ctx context.Context, telegramID int64) (*model.Form, error)
	GetLastForm(ctx context.Context, telegramID int64) (*model.Form, error)
	GetAllUserForms(ctx context.Context, telegramID int64) ([]*model.Form, error)
	GetAllForms(ctx context.Context) ([]*model.Form, error)
}

var Models = []interface{}{model.User{}, model.Token{}, model.Form{}}

type database struct {
	db     *gorm.DB
	logger *zap.Logger

	uuidGen func() uuid.UUID
}

var _ Database = database{}

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

func (d database) CreateOrProlongToken(ctx context.Context, telegramID int64) (*model.Token, error) {
	q := query.Use(d.db)
	t := q.Token

	uid := d.uuidGen().String()
	token := &model.Token{
		UUID:           uid,
		UserTelegramId: telegramID,
		ExpireAt:       time.Now().Add(time.Hour * 24),
	}
	err := t.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_telegram_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"expire_at"}),
		},
	).Create(token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (d database) ReissueToken(ctx context.Context, telegramID int64) (*model.Token, error) {
	var token *model.Token
	q := query.Use(d.db)
	t := q.Token
	err := q.Transaction(func(tx *query.Query) error {
		uid := d.uuidGen().String()
		_, err := tx.Token.WithContext(ctx).Where(t.UserTelegramId.Eq(telegramID)).Updates(&model.Token{
			UUID:     uid,
			ExpireAt: time.Now().Add(time.Hour * 24),
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

func (d database) GetUserByToken(ctx context.Context, token string) (*model.User, error) {
	q := query.Use(d.db)
	t := q.Token
	u := q.User
	user, err := u.WithContext(ctx).Join(t, u.TelegramID.EqCol(t.UserTelegramId)).Where(t.UUID.Eq(token)).Where(t.ExpireAt.GtCol(t.ExpireAt.Now())).Take()
	if err != nil {
		return nil, err
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

func (d database) UpdateUserByID(ctx context.Context, id uint, user *model.User) (*model.User, error) {
	u := query.Use(d.db).User
	userDo := u.WithContext(ctx)
	_, err := userDo.Where(u.ID.Eq(id)).Updates(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (d database) CreateForm(ctx context.Context, form *model.Form) (*model.Form, error) {
	f := query.Use(d.db).Form
	err := f.WithContext(ctx).Create(form)
	if err != nil {
		return nil, err
	}
	return form, nil
}

func (d database) AcceptForm(ctx context.Context, id uint) (*gen.ResultInfo, error) {
	f := query.Use(d.db).Form
	result, err := f.WithContext(ctx).Where(f.ID.Eq(id)).Update(f.Status, model.FormStatusAccepted)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (d database) RejectForm(ctx context.Context, id uint) (*gen.ResultInfo, error) {
	f := query.Use(d.db).Form
	result, err := f.WithContext(ctx).Where(f.ID.Eq(id)).Update(f.Status, model.FormStatusRejected)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (d database) GetActualForm(ctx context.Context, telegramID int64) (*model.Form, error) {
	f := query.Use(d.db).Form
	form, err := f.WithContext(ctx).Where(f.UserTelegramId.Eq(telegramID)).Where(f.Status.Eq(model.FormStatusAccepted)).Order(f.CreatedAt.Desc()).First()
	if err != nil {
		return nil, err
	}
	return form, nil
}

func (d database) GetLastForm(ctx context.Context, telegramID int64) (*model.Form, error) {
	f := query.Use(d.db).Form
	form, err := f.WithContext(ctx).Where(f.UserTelegramId.Eq(telegramID)).Order(f.CreatedAt.Desc()).First()
	if err != nil {
		return nil, err
	}
	return form, nil
}

func (d database) GetAllUserForms(ctx context.Context, telegramID int64) ([]*model.Form, error) {
	f := query.Use(d.db).Form
	all, err := f.WithContext(ctx).Where(f.UserTelegramId.Eq(telegramID)).Find()
	if err != nil {
		return nil, err
	}
	return all, nil
}

func (d database) GetAllForms(ctx context.Context) ([]*model.Form, error) {
	f := query.Use(d.db).Form
	all, err := f.WithContext(ctx).Find()
	if err != nil {
		return nil, err
	}
	return all, nil
}

func NewDatabase(dsn string, logger *zap.Logger) (Database, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return NewDatabaseWithDb(db, logger), nil
}

func NewDatabaseWithDb(db *gorm.DB, logger *zap.Logger) Database {
	return &database{db, logger, uuid.New}
}
