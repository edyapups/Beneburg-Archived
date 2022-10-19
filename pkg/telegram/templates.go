package telegram

import (
	"beneburg/pkg/database/model"
	"fmt"
	"strings"
)

const (
	userTelegramIDDescription  = "Telegram ID"
	userUsernameDescription    = "Username"
	userNameDescription        = "Имя"
	userAgeDescription         = "Возраст"
	userSexDescription         = "Пол"
	userAboutDescription       = "О себе"
	userHobbiesDescription     = "Хобби"
	userWorkDescription        = "Работа"
	userEducationDescription   = "Образование"
	userCoverLetterDescription = "Почему хочет к нам?"
	userContactsDescription    = "Контакты"
)

type Templator interface {
	InfoCommandNoReply() string
	InfoCommandNoUser() string
	InfoCommandReply(user *model.User) string
}

type templator struct {
}

func (t templator) InfoCommandNoUser() string {
	return "У меня нет информации об этом пользователе"
}

func (t templator) InfoCommandReply(user *model.User) string {
	stringBuilder := strings.Builder{}
	stringBuilder.WriteString(FieldStringWithBoldKey(userNameDescription, user.Name))
	AddDelimiter(&stringBuilder)

	if user.Age != nil {
		stringBuilder.WriteString(FieldStringWithBoldKey(userAgeDescription, *user.Age))
		AddDelimiter(&stringBuilder)
	}

	if user.Sex != nil {
		stringBuilder.WriteString(FieldStringWithBoldKey(userSexDescription, *user.Sex))
		AddDelimiter(&stringBuilder)
	}

	if user.About != nil {
		stringBuilder.WriteString(FieldStringWithBoldKey(userAboutDescription, *user.About))
		AddDelimiter(&stringBuilder)
	}

	if user.Hobbies != nil {
		stringBuilder.WriteString(FieldStringWithBoldKey(userHobbiesDescription, *user.Hobbies))
		AddDelimiter(&stringBuilder)
	}

	if user.Work != nil {
		stringBuilder.WriteString(FieldStringWithBoldKey(userWorkDescription, *user.Work))
		AddDelimiter(&stringBuilder)
	}

	if user.Education != nil {
		stringBuilder.WriteString(FieldStringWithBoldKey(userEducationDescription, *user.Education))
		AddDelimiter(&stringBuilder)
	}

	if user.CoverLetter != nil {
		stringBuilder.WriteString(FieldStringWithBoldKey(userCoverLetterDescription, *user.CoverLetter))
		AddDelimiter(&stringBuilder)
	}

	if user.Contacts != nil {
		stringBuilder.WriteString(FieldStringWithBoldKey(userContactsDescription, *user.Contacts))
		AddDelimiter(&stringBuilder)
	}

	stringBuilder.WriteString(fmt.Sprintf(
		"<i>ID <a href=\"tg://user?id=%d\">пользователя</a>: </i><code>%d</code>",
		user.TelegramID,
		user.TelegramID,
	))

	return stringBuilder.String()
}

func AddDelimiter(stringBuilder *strings.Builder) {
	stringBuilder.WriteString("\n\n")
}

func FieldStringWithBoldKey[T any](key string, value T) string {
	return fmt.Sprintf("%s:\n%v", FormatBold(key), value)
}

func (t templator) InfoCommandNoReply() string {
	return "Для получения информации об участнике необходимо ответить на его сообщение командой /info"
}

func NewTemplator() Templator {
	return &templator{}
}

func FormatBold(text string) string {
	return "<b>" + text + "</b>"
}

func FormatItalic(text string) string {
	return "<i>" + text + "</i>"
}
