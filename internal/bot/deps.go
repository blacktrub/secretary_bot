package bot

import (
	"context"
	"secretary_bot/internal/dto"
)

type repository interface {
	SaveMessage(ctx context.Context, msg dto.Message) error
}

type aliaser interface {
	Alias(userID int64) string
}
