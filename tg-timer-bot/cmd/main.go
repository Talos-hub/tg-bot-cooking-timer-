package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	conf "github.com/Talos-hub/tg-bot-cooking-timer-/pkg/confloader"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func setupLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	return logger

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

	// create newBotApi
	bot, err := tg.NewBotAPI(c.Token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message == nil {
			continue
		}

		msg := tg.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case "help":
			msg.Text = `Это просто обычный бот таймер, для приготовления яйца или мяса.
            Вы можете выбрать команду, и идти на ху№, ибо ты че тупой, нах!й я вобще это расписываю.
            `
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}
		case "set":
			msg.Text = `Вы можете устоновить пользователские настройки таймера.
			выберите тип еды, Мясо: /meat, Яйцо: /egg `
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}

		case "meat":
			var data string = `"meat":{
		    "second": 0,
		    "minute": 20,
 		    "hours": 0
            },`
			msg.Text = "Пришлите данные форматом: " + data
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}

		default:
			msg.Text = "На заморском я не понимаю, пиши нормально шебень"
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}
		}
	}

}
