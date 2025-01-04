package repository

import (
	"context"
	"fmt"
	"secretary_bot/internal/dto"

	"github.com/jmoiron/sqlx"
)

const schema = `
    create table if not exists message (
        id integer primary key autoincrement,
        user_id integer,
        chat_id integer,
        name text,
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
        (user_id, chat_id, name, message, created_at) 
        values (:user_id, :chat_id, :name, :message, :created_at)
    `

	if _, err := r.db.NamedExecContext(ctx, query, msg); err != nil {
		return fmt.Errorf("insert: %w", err)
	}

	return nil
}
