package database

import (
	"beneburg/pkg/database/model"
	"beneburg/pkg/utils"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
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
	t.Run("CreateToken", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `tokens` (`user_telegram_id`,`expire_at`,`uuid`) VALUES (?,?,?)")).WithArgs(
			10,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tokens` WHERE `tokens`.`uuid` = ? ORDER BY `tokens`.`uuid` LIMIT 1")).WithArgs().WillReturnRows(sqlmock.NewRows([]string{"uuid", "user_telegram_id", "expire_at"}).AddRow("test", "10", time.Now()))
		mock.ExpectCommit()
		token, err := db.CreateToken(ctx, 10)
		assert.NoError(t, err)
		assert.Equal(t, token.UserTelegramId, int64(10))
		assert.Equal(t, token.UUID, "test")
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
			Username:    utils.GetAddress("test"),
			Name:        "test",
			Age:         utils.GetAddress(int32(10)),
			Sex:         "test",
			About:       utils.GetAddress("test"),
			Hobbies:     utils.GetAddress("test"),
			Work:        utils.GetAddress("test"),
			Education:   utils.GetAddress("test"),
			CoverLetter: utils.GetAddress("test"),
			Contacts:    utils.GetAddress("test"),
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
	t.Run("UpdateUserByID - empty", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `created_at`=?, `updated_at`=?,`deleted_at`=?,`telegram_id`=?,`username`=?,`name`=?,`age`=?,`sex`=?,`about`=?,`hobbies`=?,`work`=?,`education`=?,`cover_letter`=?,`contacts`=?,`is_bot`=?,`is_active`=? WHERE `id` = ?")).
			WithArgs().WillReturnResult(sqlmock.NewResult(10, 1))
		mock.ExpectCommit()
		out, err := db.UpdateUserByID(ctx, 10, &model.User{})
		assert.NoError(t, err)
		assert.Equal(t, 1, out.RowsAffected)
	})
	t.Run("UpdateUserByID - selected", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `updated_at`=?,`username`=? WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL")).
			WithArgs(
				sqlmock.AnyArg(),
				"test",
				10,
			).WillReturnResult(sqlmock.NewResult(10, 1))
		mock.ExpectCommit()
		out, err := db.UpdateUserByID(ctx, 10, &model.User{
			Username: utils.GetAddress("test"),
			Name:     "gwegaw",
			Age:      utils.GetAddress(int32(20)),
			Sex:      "default",
			About:    utils.GetAddress("test"),
		}, "username")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), out.RowsAffected)
	})
	t.Run("UpdateOrCreateUser", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`telegram_id`,`username`,`name`,`age`,`sex`,`about`,`hobbies`,`work`,`education`,`cover_letter`,`contacts`,`is_bot`,`is_active`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `username`=VALUES(`username`)")).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				10,
				"test",
				"",
				nil,
				"undefined",
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				false,
				true,
			).WillReturnResult(sqlmock.NewResult(10, 1))
		mock.ExpectCommit()
		out, err := db.UpdateOrCreateUser(ctx, &model.User{
			TelegramID: 10,
			Username:   utils.GetAddress("test"),
			IsActive:   true,
		})
		assert.NoError(t, err)
		assert.Equal(t, int64(10), out.TelegramID)
	})
}
