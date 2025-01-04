package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"secretary_bot/internal/bot"
	"secretary_bot/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx, cn := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cn()

	db, err := sqlx.Open("sqlite3", "bot.db")
	if err != nil {
		log.Fatalf("open db: %s", err)
	}

	repo := repository.New(db)
	if err := repo.Init(ctx); err != nil {
		log.Fatalf("init schema: %s", err)
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("load env: %s", err)
	}

	bot, err := bot.New(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("init bot: %s", err)
	}

	if err := bot.Run(ctx); err != nil {
		log.Fatalf("run bot: %s", err)
	}
}
