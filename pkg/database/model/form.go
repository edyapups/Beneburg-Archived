package model

import "gorm.io/gorm"

const TableNameForm = "forms"

const (
	FormStatusNew      = "new"
	FormStatusAccepted = "accepted"
	FormStatusRejected = "rejected"
)

type Form struct {
	gorm.Model
	UserTelegramId int64 `gorm:"column:user_telegram_id" json:"user_telegram_id"`
	User           User  `gorm:"foreignKey:UserTelegramId;references:TelegramID" json:"user"`

	Name        string  `gorm:"column:name" json:"name"`
	Age         *int32  `gorm:"column:age" json:"age"`
	Gender      string  `gorm:"column:gender; type:enum('male', 'female', 'nonbinary', 'undefined');default:'undefined'" json:"gender"`
	About       *string `gorm:"column:about" json:"about"`
	Hobbies     *string `gorm:"column:hobbies" json:"hobbies"`
	Work        *string `gorm:"column:work" json:"work"`
	Education   *string `gorm:"column:education" json:"education"`
	CoverLetter *string `gorm:"column:cover_letter" json:"cover_letter"`
	Contacts    *string `gorm:"column:contacts" json:"contacts"`

	Status string `gorm:"column:status; type:enum('new', 'accepted', 'rejected');default:'new'" json:"status"`
}

func (u *Form) RuGender() string {
	if u == nil {
		return ""
	}
	switch u.Gender {
	case "male":
		return "мужской"
	case "female":
		return "женский"
	case "nonbinary":
		return "небинарный"
	default:
		return "не указан"
	}
}

const (
	UserNameDescription        = "Имя"
	UserAgeDescription         = "Возраст"
	UserGenderDescription      = "Пол"
	UserAboutDescription       = "О себе"
	UserHobbiesDescription     = "Хобби"
	UserWorkDescription        = "Работа"
	UserEducationDescription   = "Образование"
	UserCoverLetterDescription = "Почему хочет к нам?"
	UserContactsDescription    = "Контакты"
)

func (*Form) TableName() string {
	return TableNameForm
}
