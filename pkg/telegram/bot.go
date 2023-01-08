package telegram

import (
	"beneburg/pkg/database"
	"beneburg/pkg/database/model"
	"beneburg/pkg/utils"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"strings"
	"time"
)

//go:generate mockgen -source=bot.go -destination=./mocks/mock_bot.go -package=mock_telegram
type Bot struct {
	bot       TgBotAPI
	db        database.Database
	templator Templator

	messageQueue chan tgbotapi.Update

	ctx    context.Context
	logger *zap.Logger
}

func NewBot(ctx context.Context, bot TgBotAPI, db database.Database, logger *zap.Logger) *Bot {
	return &Bot{
		bot:       bot,
		templator: NewTemplator(),
		db:        db,
		ctx:       ctx,
		logger:    logger,
	}
}

func (b *Bot) Start() {
	b.messageQueue = make(chan tgbotapi.Update, 100)
	go b.startGettingUpdates()
	go b.startProcessingUpdates()
}

func (b *Bot) startGettingUpdates() {
	var offset = 0
	for {
		select {
		case <-b.ctx.Done():
			return
		default:
		}
		updates, err := b.bot.GetUpdates(tgbotapi.UpdateConfig{
			Offset:  offset,
			Timeout: 0,
		})
		if err != nil {
			b.logger.Error("Error while getting bot updates", zap.Error(err))
			b.logger.Info("Sleeping for 3 seconds...")
			time.Sleep(time.Second * 3)
			continue
		}
		for _, update := range updates {
			if update.UpdateID >= offset {
				offset = update.UpdateID + 1
				b.messageQueue <- update
			}
		}
	}
}

func (b *Bot) startProcessingUpdates() {
	for {
		select {
		case <-b.ctx.Done():
			return
		case update := <-b.messageQueue:
			b.processUpdate(update)
		}
	}
}

func (b *Bot) processUpdate(update tgbotapi.Update) {
	b.logger.Named("processUpdate").Info("Processing update", zap.Int("update_id", update.UpdateID))
	b.logger.Named("processUpdate").Debug("Update", zap.Any("update", update))
	if update.Message != nil {
		b.processMessage(update.Message)
	}
}

func (b *Bot) processPrivateMessage(message *tgbotapi.Message) {
	b.logger.Named("processPrivateMessage").Debug("Processing private message", zap.Any("message", message))
	if message.IsCommand() {
		b.processPrivateCommand(message)
		return
	}

	if message.Text == "ping" {
		b.processPing(message)
		return
	}
}

func (b *Bot) processGroupMessage(message *tgbotapi.Message) {
	if message.IsCommand() {
		b.processGroupCommand(message)
		return
	}

	if message.Text == "ping" {
		b.processPing(message)
		return
	}
}

func (b *Bot) processGroupCommand(message *tgbotapi.Message) {
	if message.Command() == "info" {
		b.processInfoCommand(message)
		return
	}
}

func (b *Bot) processPing(message *tgbotapi.Message) {
	b.logger.Named("processPing").Debug("Processing ping", zap.Any("message", message))
	_, err := b.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "pong"))
	if err != nil {
		b.logger.Named("processPing").Error("Error while sending message", zap.Error(err))
		return
	}
}

func (b *Bot) processPrivateCommand(message *tgbotapi.Message) {
	b.logger.Named("processPrivateCommand").Debug("Processing private command", zap.Any("message", message))
	if message.Command() == "login" {
		b.processLoginCommand(message)
		return
	}
}

func (b *Bot) processInfoCommand(message *tgbotapi.Message) {
	if message.ReplyToMessage == nil || message.ReplyToMessage.From == nil {
		_, err := b.bot.Send(tgbotapi.NewMessage(message.Chat.ID, b.templator.InfoCommandNoReply()))
		if err != nil {
			b.logger.Named("processInfoCommand").Error("Error while sending message", zap.Error(err))
			return
		}
		return
	}
	user, err := b.db.GetUserByTelegramID(b.ctx, message.ReplyToMessage.From.ID)
	if err != nil {
		b.logger.Named("processInfoCommand").Error("Error while getting user from db", zap.Error(err))
		return
	}
	if user == nil {
		_, err := b.bot.Send(tgbotapi.NewMessage(message.Chat.ID, b.templator.InfoCommandNoUser()))
		if err != nil {
			b.logger.Named("processInfoCommand").Error("Error while sending message", zap.Error(err))
			return
		}
		return
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, b.templator.InfoCommandReply(user))
	msg.ParseMode = tgbotapi.ModeHTML
	_, err = b.bot.Send(msg)
	if err != nil {
		b.logger.Named("processInfoCommand").Error("Error while sending message", zap.Error(err))
		return
	}
}

func (b *Bot) processMessage(message *tgbotapi.Message) {
	if from := message.From; from != nil && !from.IsBot {
		user := model.User{
			TelegramID: from.ID,
			Username: func() *string {
				if from.UserName != "" {
					return utils.GetAddress(from.UserName)
				} else {
					return nil
				}
			}(),
			Name:     strings.TrimSpace(from.FirstName + " " + from.LastName),
			IsActive: true,
		}
		_, err := b.db.UpdateOrCreateUser(b.ctx, &user)
		if err != nil {
			b.logger.Named("processMessage").Error("Error while updating user", zap.Error(err))
			return
		}
	}
	if message.Chat != nil && message.Chat.Type == "private" {
		b.processPrivateMessage(message)
	}

	if message.Chat != nil && (message.Chat.Type == "group" || message.Chat.Type == "supergroup") {
		b.processGroupMessage(message)
	}
}

func (b *Bot) processLoginCommand(message *tgbotapi.Message) {
	b.logger.Named("processLoginCommand").Debug("Processing login command", zap.Any("message", message))
	if message.From == nil {
		b.logger.Named("processLoginCommand").Error("Message's From is nil")
		return
	}
	token, err := b.db.CreateToken(b.ctx, message.From.ID)
	if err != nil {
		b.logger.Named("processLoginCommand").Error("Error while creating token", zap.Error(err))
		return
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, b.templator.LoginCommandReply(token))
	msg.ParseMode = tgbotapi.ModeHTML
	_, err = b.bot.Send(msg)
	if err != nil {
		b.logger.Named("processLoginCommand").Error("Error while sending message", zap.Error(err))
		return
	}
}
