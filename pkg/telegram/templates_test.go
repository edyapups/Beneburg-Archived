package telegram

import (
	"beneburg/pkg/database/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddField(t *testing.T) {
	result := FieldStringWithBoldKey("key", 1)
	assert.Equal(t, result, "<b>key</b>: 1")
}

func Test_templator_InfoCommandReply(t1 *testing.T) {
	result := templator{}.InfoCommandReply(&model.User{
		TelegramID:  10,
		Username:    getAddress("username"),
		Name:        "Name",
		Age:         getAddress(int32(10)),
		Sex:         getAddress("sex"),
		About:       getAddress("about"),
		Hobbies:     getAddress("hobbies"),
		Education:   getAddress("education"),
		CoverLetter: getAddress("cover_letter"),
		Contacts:    getAddress("contacts"),
	})
	assert.Equal(t1, "<b>Имя</b>:\nName\n\n"+
		"<b>Возраст</b>:\n10\n\n"+
		"<b>Пол</b>:\nsex\n\n"+
		"<b>О себе</b>:\nabout\n\n"+
		"<b>Хобби</b>:\nhobbies\n\n"+
		"<b>Образование</b>:\neducation\n\n"+
		"<b>Почему хочет к нам?</b>:\ncover_letter\n\n"+
		"<b>Контакты</b>:\ncontacts\n\n"+
		"<i>ID <a href=\"tg://user?id=10\">пользователя</a>: </i><code>10</code>", result)
}

func getAddress[T any](s T) *T {
	return &s
}
