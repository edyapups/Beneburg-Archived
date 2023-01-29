package model

import "gorm.io/gorm"

const TableNameUser = "users"

const (
	UserStatusNew       = "new"
	UserStatusActive    = "active"
	UserStatusNotActive = "not_active"
	UserStatusBot       = "bot"
	UserStatusBanned    = "banned"
)

type User struct {
	gorm.Model
	TelegramID int64   `gorm:"column:telegram_id;primaryKey;uniqueIndex:telegram_id,priority:1" json:"telegram_id"`
	Username   *string `gorm:"column:username" json:"username"`

	Status string `gorm:"column:status; type:enum('new', 'active', 'not_active', 'bot', 'banned'); default:'new'" json:"status"`
}

func (*User) TableName() string {
	return TableNameUser
}

const (
	UserTelegramIDDescription = "Telegram ID"
	UserUsernameDescription   = "Username"
)
