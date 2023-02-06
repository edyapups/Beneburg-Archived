package model

import "gorm.io/gorm"

const TableNameUser = "users"

const (
	UserStatusNew       = "new"
	UserStatusActive    = "active"
	UserStatusNotActive = "not_active"
	UserStatusAccepted  = "accepted"
	UserStatusRejected  = "rejected"
	UserStatusBot       = "bot"
	UserStatusBanned    = "banned"
)

type User struct {
	gorm.Model
	TelegramID int64   `gorm:"column:telegram_id;primaryKey;uniqueIndex:telegram_id,priority:1" json:"telegram_id"`
	Username   *string `gorm:"column:username" json:"username"`
	FirstName  string  `gorm:"column:first_name; default:''" json:"first_name"`
	LastName   *string `gorm:"column:last_name" json:"last_name"`

	Status string `gorm:"column:status; type:enum('new', 'active', 'not_active', 'accepted', 'rejected', 'bot', 'banned'); default:'new'" json:"status"`
}

func (*User) TableName() string {
	return TableNameUser
}

const (
	UserTelegramIDDescription = "Telegram ID"
	UserUsernameDescription   = "Username"
)
