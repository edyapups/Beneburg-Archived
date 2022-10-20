package database

import (
	"beneburg/pkg/database/model"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

func Test_database_AutoMigrate(t *testing.T) {
	ctx := context.Background()
	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)

	mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("8.0.0"))

	engine, err := gorm.Open(mysql.New(mysql.Config{Conn: dbMock}), &gorm.Config{})
	assert.NoError(t, err)
	db := NewDatabaseWithDb(engine, nil)

	t.Run("Get user by id", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
			WithArgs(10).
			WillReturnRows(sqlmock.NewRows([]string{
				"id",
				"telegram_id",
				"username",
				"name",
				"age",
				"sex",
				"about",
				"hobbies",
				"work",
				"education",
				"cover_letter",
				"contacts",
				"is_bot",
				"is_active"}).AddRow(
				10, 0, "", "", 0, "", "", "", "", "", "", "", false, false))
		user, err := db.GetUserByID(ctx, 10)
		assert.NoError(t, err)
		assert.Equal(t, uint(10), user.ID)
	})

	t.Run("CreateUser", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`telegram_id`,`username`,`name`,`age`,`sex`,`about`,`hobbies`,`work`,`education`,`cover_letter`,`contacts`,`is_bot`,`is_active`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				10,
				"test",
				"test",
				10,
				"test",
				"test",
				"test",
				"test",
				"test",
				"test",
				"test",
				false,
				true,
			).WillReturnResult(sqlmock.NewResult(10, 1))
		mock.ExpectCommit()

		_, err := db.CreateUser(ctx, &model.User{
			TelegramID:  10,
			Username:    getAddress("test"),
			Name:        "test",
			Age:         getAddress(int32(10)),
			Sex:         "test",
			About:       getAddress("test"),
			Hobbies:     getAddress("test"),
			Work:        getAddress("test"),
			Education:   getAddress("test"),
			CoverLetter: getAddress("test"),
			Contacts:    getAddress("test"),
			IsBot:       false,
			IsActive:    true,
		})
		assert.NoError(t, err)
	})
	t.Run("GetUserIDByToken", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT `users`.`id` FROM `users` INNER JOIN `tokens` ON `users`.`telegram_id` = `tokens`.`user_telegram_id` WHERE `tokens`.`uuid` = ? AND `tokens`.`expire_at` > ? AND `users`.`deleted_at` IS NULL LIMIT 1")).WithArgs("test", sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))
		out, err := db.GetUserIDByToken(ctx, "test")
		assert.Equal(t, uint(10), out)
		assert.NoError(t, err)
	})
	t.Run("GetUserIDByToken2", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT `users`.`id` FROM `users` INNER JOIN `tokens` ON `users`.`telegram_id` = `tokens`.`user_telegram_id` WHERE `tokens`.`uuid` = ? AND `tokens`.`expire_at` > ? AND `users`.`deleted_at` IS NULL LIMIT 1")).WithArgs("test", sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{}))
		_, err := db.GetUserIDByToken(ctx, "test")
		assert.Error(t, err)
	})
}

func getAddress[T any](s T) *T {
	return &s
}
