package repository

import (
	"context"
	"fmt"
	"secretary_bot/internal/dto"
	"time"

	"github.com/jmoiron/sqlx"
)

const schema = `
    create table if not exists message (
        id integer primary key autoincrement,
        reporter_id integer,
        reporter_name string,
        user_id integer,
        user_name text,
        chat_id integer,
        message text,
        created_at integer
    )
`

type repo struct {
	db *sqlx.DB
}

func New(
	db *sqlx.DB,
) *repo {
	return &repo{
		db: db,
	}
}

func (r *repo) Init(ctx context.Context) error {
	if _, err := r.db.ExecContext(ctx, schema); err != nil {
		return fmt.Errorf("init schema: %w", err)
	}

	return nil
}

func (r *repo) SaveMessage(ctx context.Context, msg dto.Message) error {
	const query = `
        insert into message 
        (
            reporter_id,
            reporter_name,
            user_id,
            user_name,
            chat_id,
            message,
            created_at
        ) 
        values (
            :reporter_id,
            :reporter_name,
            :user_id,
            :user_name,
            :chat_id,
            :message,
            :created_at
        )
    `

	if _, err := r.db.NamedExecContext(ctx, query, msg); err != nil {
		return fmt.Errorf("insert: %w", err)
	}

	return nil
}

func (r *repo) Retro(ctx context.Context, chatID int64, start time.Time) ([]dto.Message, error) {
	const query = `
        select reporter_id, reporter_name, user_id, user_name, chat_id, message, created_at
        from message 
        where chat_id = $1 and created_at > $2
	`

	rows := make([]dto.Message, 0, 100)
	if err := r.db.SelectContext(ctx, &rows, query, chatID, start.Unix()); err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}

	return rows, nil
}
