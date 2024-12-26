package bot

import (
	"context"
	"fmt"
	"regexp"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	timeout   = 10
	ownerName = "Дмитрий Евгеньевич"
)

var (
	msgRegex = regexp.MustCompile(`.*([на|На]\sретро).*`)
)

type bot struct {
	api     *api.BotAPI
	ownerID int64
}

func New(token string, ownerID int64) (*bot, error) {
	b, err := api.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("init bot: %w", err)
	}

	return &bot{
		api:     b,
		ownerID: ownerID,
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
				reply := "Записала!"
				if msg.From.ID == b.ownerID {
					reply = fmt.Sprintf("%s, записала!", ownerName)

				}

				b.reply(msg.Chat.ID, msg.MessageID, reply)
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
