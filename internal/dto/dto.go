package dto

type Message struct {
	ReporterID   int64  `db:"reporter_id"`
	ReporterName string `db:"reporter_name"`
	UserID       int64  `db:"user_id"`
	UserName     string `db:"user_name"`
	ChatID       int64  `db:"chat_id"`
	Message      string `db:"message"`
	CreatedAt    int64  `db:"created_at"`
}
