package bot

import (
	"context"
	"fmt"
	"regexp"
	"secretary_bot/internal/dto"
	"time"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	timeout = 10
)

var (
	msgRegex = regexp.MustCompile(`.*([на|На]\sретро).*`)
)

type bot struct {
	api           *api.BotAPI
	repo          repository
	answerBuilder answerBuilder
}

func New(
	token string,
	repo repository,
	aliaser aliaser,
) (*bot, error) {
	b, err := api.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("init bot: %w", err)
	}

	return &bot{
		api:  b,
		repo: repo,
		answerBuilder: answerBuilder{
			aliaser: aliaser,
		},
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
				reply := msg.ReplyToMessage
				err := b.repo.SaveMessage(ctx, dto.Message{
					ReporterID:   msg.From.ID,
					ReporterName: msg.From.String(),
					UserID:       reply.From.ID,
					UserName:     reply.From.String(),
					ChatID:       msg.Chat.ID,
					Message:      reply.Text,
					CreatedAt:    time.Now(),
				})
				if err != nil {
					fmt.Printf("save err: %s\n", err)
					continue
				}

				b.reply(msg.Chat.ID, msg.MessageID, b.answerBuilder.Msg(msg.From.ID))
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
