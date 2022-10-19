package telegram

import (
	"beneburg/pkg/database"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"time"
)

//go:generate mockgen -source=bot.go -destination=./mocks/mock_bot.go -package=mock_telegram
type Bot struct {
	bot TgBotAPI
	db  database.Database

	messageQueue chan tgbotapi.Update

	ctx    context.Context
	logger *zap.Logger
}

func NewBot(ctx context.Context, bot TgBotAPI, db database.Database, logger *zap.Logger) *Bot {
	return &Bot{
		bot:    bot,
		db:     db,
		ctx:    ctx,
		logger: logger,
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
	if update.Message != nil && update.Message.Chat != nil && update.Message.Chat.Type == "private" {
		b.logger.Debug("New private message",
			zap.String("user", update.Message.From.String()),
			zap.String("text", update.Message.Text),
		)
		b.processPrivateMessage(update.Message)
	}

	if update.Message != nil && update.Message.Chat != nil && update.Message.Chat.Type == "group" {
		b.logger.Debug("New group message",
			zap.String("user", update.Message.From.String()),
			zap.String("chat", update.Message.Chat.Title),
			zap.String("text", update.Message.Text),
		)
		b.processGroupMessage(update.Message)
	}
}

func (b *Bot) processPrivateMessage(message *tgbotapi.Message) {
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
	b.logger.Named("processGroupCommand").Warn("Not implemented")
}

func (b *Bot) processPing(message *tgbotapi.Message) {
	_, err := b.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "pong"))
	if err != nil {
		b.logger.Named("processPing").Error("Error while sending message", zap.Error(err))
		return
	}
}

func (b *Bot) processPrivateCommand(message *tgbotapi.Message) {
	b.logger.Named("processPrivateCommand").Warn("Not implemented")
}
