package bot

import (
	"context"
	"secretary_bot/internal/dto"
	"time"
)

type repository interface {
	SaveMessage(ctx context.Context, msg dto.Message) error
	Retro(ctx context.Context, chatID int64, start time.Time) ([]dto.Message, error)
}

type aliaser interface {
	Alias(userID int64) string
}
