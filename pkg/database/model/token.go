package model

import "time"

const TableNameToken = "tokens"

type Token struct {
	UUID           string    `gorm:"column:uuid;primaryKey;uniqueIndex:uuid,priority:1" json:"uuid"`
	UserTelegramId int64     `gorm:"column:user_telegram_id" json:"user_telegram_id"`
	User           User      `gorm:"foreignKey:UserTelegramId;references:TelegramID" json:"user"`
	ExpireAt       time.Time `gorm:"column:expire_at" json:"expire_at"`
}

func (*Token) TableName() string {
	return TableNameToken
}
