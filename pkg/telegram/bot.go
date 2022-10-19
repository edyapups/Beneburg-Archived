package telegram

import (
	"context"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"time"
)

//go:generate mockgen -source=bot.go -destination=./mocks/mock_bot.go -package=mock_telegram

type Bot struct {
	bot TgBotAPI

	messageQueue chan interface{}
	rateLimiter  *rate.Limiter

	ctx    context.Context
	logger *zap.Logger
}

func NewBot(ctx context.Context, logger *zap.Logger, bot TgBotAPI) *Bot {
	return &Bot{
		bot:    bot,
		ctx:    ctx,
		logger: logger,
	}
}

func (b *Bot) Start() {
	b.messageQueue = make(chan interface{})
	b.rateLimiter = rate.NewLimiter(rate.Every(time.Second), 5)
	go b.processMessages()
}

func (b *Bot) processMessages() {
	for {
		select {
		case <-b.ctx.Done():
			return
		case msg := <-b.messageQueue:
			err := b.rateLimiter.Wait(b.ctx)
			if err != nil {
				b.logger.Error("failed to wait for rate limiter", zap.Error(err))
				return
			}
			b.processMessage(msg)
		}
	}
}

func (b *Bot) processMessage(msg interface{}) {
	//TODO implement me
	panic("implement me")
}
