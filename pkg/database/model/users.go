package model

import "gorm.io/gorm"

const TableNameUser = "users"

type User struct {
	gorm.Model
	TelegramID  int32   `gorm:"column:telegram_id;;primaryKey;uniqueIndex:telegram_id,priority:1" json:"telegram_id"`
	Username    *string `gorm:"column:username;" json:"username"`
	Name        string  `gorm:"column:name;;not null" json:"name"`
	Age         *int32  `gorm:"column:age;" json:"age"`
	Sex         *string `gorm:"column:sex;" json:"sex"`
	About       *string `gorm:"column:about;" json:"about"`
	Hobbies     *string `gorm:"column:hobbies;" json:"hobbies"`
	Work        *string `gorm:"column:work;" json:"work"`
	Education   *string `gorm:"column:education;" json:"education"`
	CoverLetter *string `gorm:"column:cover_letter;" json:"cover_letter"`
	Contacts    *string `gorm:"column:contacts;" json:"contacts"`
	IsBot       bool    `gorm:"column:is_bot;not null" json:"is_bot"`
	IsActive    bool    `gorm:"column:is_active;not null;default:1" json:"is_active"`
}

func (*User) TableName() string {
	return TableNameUser
}
