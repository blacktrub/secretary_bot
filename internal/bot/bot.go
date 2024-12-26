package bot

import (
	"context"
	"fmt"
	"regexp"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	timeout = 10
)

var (
	msgRegex = regexp.MustCompile(`.*([на|На]\sретро).*`)
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
			msg := update.Message
			if msg != nil && msg.ReplyToMessage != nil && msgRegex.Match([]byte(msg.Text)) {
				ans := newAnswer(msg.From.ID)
				b.reply(msg.Chat.ID, msg.MessageID, ans.Msg())
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
