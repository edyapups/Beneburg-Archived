package telegram

import (
	"beneburg/pkg/database"
	"beneburg/pkg/database/model"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"time"
)

//go:generate mockgen -source=bot.go -destination=./mocks/mock_bot.go -package=mock_telegram
type Bot interface {
	Start()
	GetSendFunc() TelegramBotSendFunc
	SetLogger(logger *zap.Logger)
}

type TelegramBotSendFunc func(message tgbotapi.Chattable)

type botManager struct {
	bot       TgBotAPI
	db        database.Database
	templator Templator

	updatesChan  chan tgbotapi.Update
	messagesChan chan tgbotapi.Chattable

	ctx    context.Context
	logger *zap.Logger
}

func NewBot(ctx context.Context, bot TgBotAPI, db database.Database) Bot {
	return &botManager{
		bot:       bot,
		templator: NewTemplator(),
		db:        db,
		ctx:       ctx,
	}
}

func (b *botManager) SetLogger(logger *zap.Logger) {
	b.logger = logger
}

// TODO: will be parallelized
func (b *botManager) GetSendFunc() TelegramBotSendFunc {
	return b.send
}

func (b *botManager) Start() {
	b.updatesChan = make(chan tgbotapi.Update, 60)
	b.messagesChan = make(chan tgbotapi.Chattable, 60)
	go b.startGettingUpdates()
	go b.startProcessingUpdates()
	go b.startProcessingMessages()
}

// TODO: add With() with context to all loggings
func (b *botManager) startGettingUpdates() {
	var offset = 0
	for {
		select {
		case <-b.ctx.Done():
			return
		default:
		}
		updates, err := b.bot.GetUpdates(tgbotapi.UpdateConfig{
			Offset:  offset,
			Timeout: 60,
		})
		b.logger.Named("startGettingUpdates").Debug("Got updates", zap.Any("updates", updates))
		if err != nil {
			b.logger.Error("Error while getting bot updates", zap.Error(err))
			b.logger.Info("Sleeping for 3 seconds...")
			time.Sleep(time.Second * 3)
			continue
		}
		for _, update := range updates {
			if update.UpdateID >= offset {
				offset = update.UpdateID + 1
				b.updatesChan <- update
			}
		}
	}
}

func (b *botManager) startProcessingUpdates() {
	for {
		select {
		case <-b.ctx.Done():
			return
		case update := <-b.updatesChan:
			b.processUpdate(update)
		}
	}
}

func (b *botManager) startProcessingMessages() {
	limiter := rate.NewLimiter(30, 30)
	for {
		select {
		case <-b.ctx.Done():
			return
		case message := <-b.messagesChan:
			_ = limiter.Wait(b.ctx)
			_, err := b.bot.Send(message)
			switch typedErr := err.(type) {
			case tgbotapi.Error:
				if typedErr.Code == 429 {
					b.logger.Named("startProcessingMessages").Info("Too many requests, sleeping for 1 second")
					b.messagesChan <- message
					time.Sleep(time.Second)
				}
			case nil:
				continue
			default:
				b.logger.Named("startProcessingMessages").Error("Error while sending message", zap.Error(err), zap.Any("message", message))
			}
		}
	}
}

func (b *botManager) send(message tgbotapi.Chattable) {
	b.messagesChan <- message
}

func (b *botManager) processUpdate(update tgbotapi.Update) {
	b.logger.Named("processUpdate").Info("Processing update", zap.Int("update_id", update.UpdateID))
	b.logger.Named("processUpdate").Debug("Update", zap.Any("update", update))
	if update.Message != nil {
		b.processMessage(update.Message)
	}
}

func (b *botManager) processMessage(message *tgbotapi.Message) {
	b.logger.Named("processMessage").Debug("Processing message", zap.Any("message", message))
	if from := message.From; from != nil && !from.IsBot {
		user := model.User{
			TelegramID: from.ID,
			Username: func() *string {
				if from.UserName != "" {
					return &from.UserName
				} else {
					return nil
				}
			}(),
		}
		_, err := b.db.UpdateOrCreateUser(b.ctx, &user)
		if err != nil {
			b.logger.Named("processMessage").Error("Error while updating user", zap.Error(err))
			return
		}
	}

	// TODO: use switch
	if message.Chat != nil && message.Chat.Type == "private" {
		b.processPrivateMessage(message)
	}

	if message.Chat != nil && (message.Chat.Type == "group" || message.Chat.Type == "supergroup") {
		b.processGroupMessage(message)
	}
}
func (b *botManager) processPrivateMessage(message *tgbotapi.Message) {
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

func (b *botManager) processGroupMessage(message *tgbotapi.Message) {
	if message.IsCommand() {
		b.processGroupCommand(message)
		return
	}

	if message.Text == "ping" {
		b.processPing(message)
		return
	}
}

func (b *botManager) processGroupCommand(message *tgbotapi.Message) {
	if message.Command() == "info" {
		b.processInfoCommand(message)
		return
	}
}

func (b *botManager) processPing(message *tgbotapi.Message) {
	b.logger.Named("processPing").Debug("Processing ping", zap.Any("message", message))
	b.send(tgbotapi.NewMessage(message.Chat.ID, "pong"))
}

func (b *botManager) processPrivateCommand(message *tgbotapi.Message) {
	b.logger.Named("processPrivateCommand").Debug("Processing private command", zap.Any("message", message))
	if message.Command() == "login" {
		b.processLoginCommand(message)
		return
	}
}

func (b *botManager) processInfoCommand(message *tgbotapi.Message) {
	if message.ReplyToMessage == nil || message.ReplyToMessage.From == nil {
		b.send(tgbotapi.NewMessage(message.Chat.ID, b.templator.InfoCommandNoReply()))
		return
	}
	user, err := b.db.GetUserByTelegramID(b.ctx, message.ReplyToMessage.From.ID)
	if err != nil {
		b.logger.Named("processInfoCommand").Error("Error while getting user from db", zap.Error(err))
		return
	}
	if user == nil {
		b.send(tgbotapi.NewMessage(message.Chat.ID, b.templator.InfoCommandNoUser()))
		return
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, b.templator.InfoCommandReply(user))
	msg.ParseMode = tgbotapi.ModeHTML
	b.send(msg)
}

func (b *botManager) processLoginCommand(message *tgbotapi.Message) {
	b.logger.Named("processLoginCommand").Debug("Processing login command", zap.Any("message", message))
	if message.From == nil {
		b.logger.Named("processLoginCommand").Error("Message's From is nil")
		return
	}
	token, err := b.db.CreateOrProlongToken(b.ctx, message.From.ID)
	if err != nil {
		b.logger.Named("processLoginCommand").Error("Error while creating token", zap.Error(err))
		return
	}
	b.logger.Named("processLoginCommand").Debug("Token created", zap.Any("token", token))
	msg := tgbotapi.NewMessage(message.Chat.ID, b.templator.LoginCommandReply(token))
	msg.ParseMode = tgbotapi.ModeHTML
	b.send(msg)
}
