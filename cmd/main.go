package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"secretary_bot/internal/bot"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("load env: %s", err)
	}

	owner, err := strconv.Atoi(os.Getenv("BOT_OWNER_ID"))
	bot, err := bot.New(os.Getenv("BOT_TOKEN"), int64(owner))
	if err != nil {
		log.Fatalf("init bot: %s", err)
	}

	ctx, cn := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cn()

	if err := bot.Run(ctx); err != nil {
		log.Fatalf("run bot: %s", err)
	}
}
