package bot

import (
	"context"
	"fmt"
	"regexp"
	"secretary_bot/internal/dto"
	"strings"
	"time"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	timeout = 10
)

var (
	msgRegex          = regexp.MustCompile(`.*([на|На]\sретро).*`)
	retroCommandRegex = regexp.MustCompile(`^\/.*retro$`)
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
					CreatedAt:    time.Now().Unix(),
				})
				if err != nil {
					fmt.Printf("save err: %s\n", err)
					continue
				}

				if err := b.reply(msg.Chat.ID, msg.MessageID, escape(b.answerBuilder.Msg(msg.From.ID))); err != nil {
					fmt.Printf("reply err: %s\n", err)
					continue
				}
			}

			if msg != nil && retroCommandRegex.Match([]byte(msg.Text)) {
				messages, err := b.repo.Retro(ctx, msg.Chat.ID, time.Now().Add(-10*24*time.Hour))
				if err != nil {
					fmt.Printf("get messages err: %s\n", err)
					continue
				}

				byReporter := make(map[int64]map[int64][]dto.Message, len(messages))
				for _, msg := range messages {
					if _, ok := byReporter[msg.ReporterID]; !ok {
						byReporter[msg.ReporterID] = map[int64][]dto.Message{}
					}

					byReporter[msg.ReporterID][msg.UserID] = append(byReporter[msg.ReporterID][msg.UserID], msg)
				}

				var report strings.Builder
				for _, messagesByReporter := range byReporter {

					var first dto.Message
					for _, byUser := range messagesByReporter {
						for _, msg := range byUser {
							first = msg
							break
						}
						break
					}

					report.WriteString(fmt.Sprintf("*%s зарепортил:*\n", first.ReporterName))

					for _, messagesByUser := range messagesByReporter {
						firstByUser := messagesByUser[0]
						report.WriteString(fmt.Sprintf("_От %s_\n", escape(firstByUser.UserName)))
						for _, msg := range messagesByUser {
							report.WriteString(fmt.Sprintf(">>%s\n\n", escape(msg.Message)))
						}
					}

					report.WriteString("\n")
				}

				b.reply(msg.Chat.ID, msg.MessageID, report.String())
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *bot) reply(chatID int64, msgID int, text string) error {
	msg := api.NewMessage(chatID, text)
	msg.ReplyToMessageID = msgID
	msg.ParseMode = api.ModeMarkdownV2
	if _, err := b.api.Send(msg); err != nil {
		return fmt.Errorf("reply: %w", err)
	}

	return nil
}

func escape(in string) string {
	return api.EscapeText(api.ModeMarkdownV2, in)
}
