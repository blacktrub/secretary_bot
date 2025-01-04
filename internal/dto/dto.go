package dto

import "time"

type Message struct {
	UserID    int64     `db:"user_id"`
	ChatID    int64     `db:"chat_id"`
	Name      string    `db:"name"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}
