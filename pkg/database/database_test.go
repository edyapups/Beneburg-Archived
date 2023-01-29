package database

import (
	"beneburg/pkg/database/model"
	"beneburg/pkg/utils"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
	"testing"
	"time"
)

func Test_database_AutoMigrate(t *testing.T) {
	ctx := context.Background()
	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)

	mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("8.0.0"))

	engine, err := gorm.Open(mysql.New(mysql.Config{Conn: dbMock}), &gorm.Config{})
	assert.NoError(t, err)
	db := NewDatabaseWithDb(engine, nil)
	var testUUID uuid.UUID
	copy(testUUID[:], []byte("test"))
	db.(*database).uuidGen = func() uuid.UUID {
		return testUUID
	}

	t.Run("Get user by id", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
			WithArgs(10).
			WillReturnRows(sqlmock.NewRows([]string{
				"id",
				"created_at",
				"updated_at",
				"deleted_at",
				"telegram_id",
				"username",
				"status",
			}).AddRow(
				10,
				time.Now(),
				time.Now(),
				nil,
				11,
				"test",
				model.UserStatusActive,
			))
		user, err := db.GetUserByID(ctx, 10)
		assert.NoError(t, err)
		assert.Equal(t, uint(10), user.ID)
	})
	t.Run("CreateOrProlongToken", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `tokens` (`user_telegram_id`,`expire_at`,`uuid`) VALUES (?,?,?)")).WithArgs(
			10,
			sqlmock.AnyArg(),
			testUUID.String(),
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		token, err := db.CreateOrProlongToken(ctx, 10)
		assert.NoError(t, err)
		assert.Equal(t, token.UserTelegramId, int64(10))
		assert.Equal(t, token.UUID, testUUID.String())
	})
	t.Run("CreateUser", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`telegram_id`,`username`,`status`) VALUES (?,?,?,?,?,?)")).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				10,
				"test",
				model.UserStatusNew,
			).WillReturnResult(sqlmock.NewResult(10, 1))
		mock.ExpectCommit()

		_, err := db.CreateUser(ctx, &model.User{
			TelegramID: 10,
			Username:   utils.GetAddress("test"),
			Status:     model.UserStatusNew,
		})
		assert.NoError(t, err)
	})

}
