package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const schema = `
    create table if not exists message (
        id integer primary key autoincrement,
        user_id integer,
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
