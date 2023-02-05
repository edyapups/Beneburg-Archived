package views

import (
	"beneburg/pkg/database"
	"beneburg/pkg/database/model"
	"beneburg/pkg/telegram"
	"beneburg/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Views interface {
	RegisterRoutes(router gin.IRouter)
	RegisterLogin(router gin.IRouter)
	RegisterProfile(router gin.IRouter)
}

var _ Views = &views{}

type views struct {
	db              database.Database
	logger          *zap.Logger
	sendToBot       telegram.TelegramBotSendFunc
	templator       telegram.Templator
	adminTelegramID int64
	groupTelegramID int64
}

func (v views) RegisterRoutes(router gin.IRouter) {
	router.GET("/", v.index)
}

func (v views) RegisterProfile(router gin.IRouter) {
	router.GET("/", v.profile)
	router.POST("/form", v.profileForm)
}

func (v views) RegisterLogin(router gin.IRouter) {
	router.GET("/", v.login)
	router.GET("/:token", v.login)
}

func NewViews(db database.Database, logger *zap.Logger, sendFunc telegram.TelegramBotSendFunc, adminTelegramID int64, groupTelegramID int64) Views {
	return &views{
		db:              db,
		logger:          logger,
		sendToBot:       sendFunc,
		templator:       telegram.NewTemplator(),
		adminTelegramID: adminTelegramID,
		groupTelegramID: groupTelegramID,
	}
}

func (v views) index(g *gin.Context) {
	_ = g.MustGet("currentUser").(*model.User)
	g.Redirect(301, "/profile")
	//g.HTML(200, "index.gohtml", gin.H{
	//	"title": "Главная",
	//	"page":  "index",
	//	"user":  user,
	//})
}

func (v views) login(g *gin.Context) {
	token := g.Param("token")
	if token == "" {
		g.Redirect(301, "https://t.me/BeneburgBot")
		return
	}
	g.SetCookie("token", token, 60*60*24, "/", "", false, true)
	g.Redirect(302, "/")
}

func (v views) profile(g *gin.Context) {
	var no_forms = false
	user := g.MustGet("currentUser").(*model.User)
	form, _ := v.db.GetLastForm(g, user.TelegramID)
	if form == nil {
		no_forms = true
		form = &model.Form{}
	}
	g.HTML(200, "profile.gohtml", gin.H{
		"title":    "Профиль",
		"page":     "profile",
		"user":     user,
		"form":     form,
		"no_forms": no_forms,
	})
}

func (v views) profileForm(g *gin.Context) {
	user := g.MustGet("currentUser").(*model.User)

	form := &model.Form{
		UserTelegramId: user.TelegramID,
	}
	nameFormValue, ok := g.GetPostForm("name")
	if ok {
		form.Name = nameFormValue
	}
	ageFormValue, ok := g.GetPostForm("age")
	if ok {
		convertedAge, err := strconv.Atoi(ageFormValue)
		if err != nil {
			v.logger.Named("profileForm").Error("Error converting age", zap.Error(err))
			g.Redirect(http.StatusFound, "/profile")
			return
		}
		form.Age = utils.GetAddress(int32(convertedAge))
	}
	genderFormValue, ok := g.GetPostForm("gender")
	if ok {
		form.Gender = genderFormValue
	}
	aboutFormValue, ok := g.GetPostForm("about")
	if ok {
		form.About = &aboutFormValue
	}
	hobbiesFormValue, ok := g.GetPostForm("hobbies")
	if ok {
		form.Hobbies = &hobbiesFormValue
	}
	workFormValue, ok := g.GetPostForm("work")
	if ok {
		form.Work = &workFormValue
	}
	educationFormValue, ok := g.GetPostForm("education")
	if ok {
		form.Education = &educationFormValue
	}
	coverLetterFormValue, ok := g.GetPostForm("cover_letter")
	if ok {
		form.CoverLetter = &coverLetterFormValue
	}
	contactsFormValue, ok := g.GetPostForm("contacts")
	if ok {
		form.Contacts = &contactsFormValue
	}
	_, err := v.db.CreateForm(g, form)
	if err != nil {
		v.logger.Named("profileForm").Error("Error creating form", zap.Error(err))
	}
	message := tgbotapi.NewMessage(user.TelegramID, v.templator.FormReceived())
	v.sendToBot(message)

	adminMessageText := fmt.Sprintf("Новая анкета:\n\n%s\n\n%s", v.templator.FormInfo(form), v.templator.UserIdWithHref(user))
	adminMessage := tgbotapi.NewMessage(v.adminTelegramID, adminMessageText)
	adminMessage.ParseMode = "HTML"
	acceptionButton := tgbotapi.NewInlineKeyboardButtonData("Принять", fmt.Sprintf("admin:form:accept:%d", form.ID))
	rejectionButton := tgbotapi.NewInlineKeyboardButtonData("Отклонить", fmt.Sprintf("admin:form:reject:%d", form.ID))
	adminMessage.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(acceptionButton, rejectionButton))
	v.sendToBot(adminMessage)

	g.Redirect(http.StatusFound, "/profile")
}
