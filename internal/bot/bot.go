package bot

import (
	"context"
	"fmt"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	timeout = 10
)

type bot struct {
	api *api.BotAPI
}

func New(token string) (*bot, error) {
	b, err := api.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("init bot: %w", err)
	}

	return &bot{
		api: b,
	}, nil
}

func (b *bot) Run(ctx context.Context) error {
	cnf := api.NewUpdate(0)
	cnf.Timeout = timeout

	ch := b.api.GetUpdatesChan(cnf)
	for {
		select {
		case update := <-ch:
			if update.Message != nil && update.Message.ReplyToMessage != nil {
				b.reply(update.Message.Chat.ID, update.Message.MessageID, "pong")
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *bot) reply(chatID int64, msgID int, text string) error {
	msg := api.NewMessage(chatID, text)
	msg.ReplyToMessageID = msgID
	if _, err := b.api.Send(msg); err != nil {
		return fmt.Errorf("reply: %w", err)
	}

	return nil
}
