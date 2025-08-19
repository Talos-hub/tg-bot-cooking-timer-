package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	bot "github.com/Talos-hub/tg-bot-cooking-timer-/pkg/bot"
	conf "github.com/Talos-hub/tg-bot-cooking-timer-/pkg/confloader"
)

func setupLogger() *slog.Logger {
	file, err := os.OpenFile("bot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// Fallback to stdout
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	// Log to file only
	return slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

}

func main() {
	// load env
	err := godotenv.Load()
	if err != nil {
		slog.Warn("env are not found")
		os.Exit(1)
	}

	logger := setupLogger()
	c, err := conf.LoadConfig(logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	b, err := bot.NewBot(c.Token, logger)
	if err != nil {
		logger.Error("error creating new bot", "error", err)
		os.Exit(1)
	}

	b.Start(&c.Food)

}
