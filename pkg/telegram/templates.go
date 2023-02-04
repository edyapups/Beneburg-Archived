package telegram

import (
	"beneburg/pkg/database/model"
	"beneburg/pkg/utils"
	"fmt"
	"html"
	"strings"
)

type Templator interface {
	InfoCommandNoReply() string
	InfoCommandNoUser() string
	InfoCommandReply(user *model.User, form *model.Form) string
	LoginCommandReply(token *model.Token) string
}

var _ Templator = templator{}

type templator struct {
}

func (t templator) LoginCommandReply(token *model.Token) string {
	stringBuilder := strings.Builder{}
	stringBuilder.WriteString("Вот ваша ссылка для входа:\n")
	stringBuilder.WriteString(html.EscapeString(utils.URLFromToken(token.UUID)))
	return stringBuilder.String()
}

func (t templator) InfoCommandNoUser() string {
	return "У меня нет информации об этом пользователе"
}

func (t templator) InfoCommandReply(user *model.User, form *model.Form) string {
	stringBuilder := strings.Builder{}
	stringBuilder.WriteString("<b>Информация об участнике:</b>")
	AddDelimiter(&stringBuilder)
	stringBuilder.WriteString(fmt.Sprintf("<b>%s</b>:\n%s", model.UserNameDescription, form.Name))
	AddDelimiter(&stringBuilder)
	if form.Age != nil {
		stringBuilder.WriteString(fmt.Sprintf("<b>%s</b>:\n%s", model.UserAgeDescription, form.Age))
		AddDelimiter(&stringBuilder)
	}
	stringBuilder.WriteString(fmt.Sprintf("<b>%s</b>:\n%s", model.UserGenderDescription, form.Gender))
	AddDelimiter(&stringBuilder)
	if form.About != nil {
		stringBuilder.WriteString(fmt.Sprintf("<b>%s</b>:\n%s", model.UserAboutDescription, form.About))
		AddDelimiter(&stringBuilder)
	}
	if form.Hobbies != nil {
		stringBuilder.WriteString(fmt.Sprintf("<b>%s</b>:\n%s", model.UserHobbiesDescription, form.Hobbies))
		AddDelimiter(&stringBuilder)
	}
	if form.Work != nil {
		stringBuilder.WriteString(fmt.Sprintf("<b>%s</b>:\n%s", model.UserWorkDescription, form.Work))
		AddDelimiter(&stringBuilder)
	}
	if form.Education != nil {
		stringBuilder.WriteString(fmt.Sprintf("<b>%s</b>:\n%s", model.UserEducationDescription, form.Education))
		AddDelimiter(&stringBuilder)
	}
	if form.CoverLetter != nil {
		stringBuilder.WriteString(fmt.Sprintf("<b>%s</b>:\n%s", model.UserCoverLetterDescription, form.CoverLetter))
		AddDelimiter(&stringBuilder)
	}
	if form.Contacts != nil {
		stringBuilder.WriteString(fmt.Sprintf("<b>%s</b>:\n%s", model.UserContactsDescription, form.Contacts))
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
