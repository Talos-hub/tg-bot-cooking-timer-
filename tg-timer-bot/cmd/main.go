package main

import (
	"fmt"
	"log"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	fmt.Println("Hello world!!!")
	bot, err := tg.NewBotAPI("______")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
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
		default:
			msg.Text = "На заморском я не понимаю, пиши нормально шебень"
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}
		}
	}

}
