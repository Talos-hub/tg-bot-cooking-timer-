package bot

import (
	"fmt"

	"github.com/Talos-hub/tg-bot-cooking-timer-/pkg/consts"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Start is launch the bot
func (b *bot) Start() {
	b.logger.Info("Authorized on account", "username", b.api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)
	// start update
	for update := range updates {
		// if it's nil then we just skip it
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		msg := tgbotapi.NewMessage(chatID, "")

		// check if we are waiting for user input
		if state, exist := b.userState[chatID]; exist && state.WaitingForInput {
			b.handleUserInput(update.Message, state)
			delete(b.userState, chatID)

			continue
		}

		// handle command
		switch update.Message.Command() {
		case consts.START, consts.HELP: // hanlde start or help command
			msg.Text = consts.TEXT_HELP
		// ------------------------------------
		case consts.SETTINGS: // setup settings
			msg.Text = consts.TEXT_SETTINGS
		// ------------------------------------
		case consts.MEAT: // setup settings for meat
			msg.Text = consts.TEXT_MEAT
			// change user state
			b.userState[chatID] = &UserState{
				WaitingForInput: true,
				FoodType:        "meat",
			}
		// --------------------------------------
		case consts.EGG: // setup settings for egg
			msg.Text = consts.TEXT_EGG
			// change user state
			b.userState[chatID] = &UserState{
				WaitingForInput: true,
				FoodType:        "egg",
			}
		// --------------------------------------------
		case consts.START_TIMER: // start timer
			// TODO
		case consts.SHOW:

			msg.Text = fmt.Sprintf("Настройки для Мяса: Часы: %d, Минуты: %d, Секунды: %d\nНастройки для Яйца: Часы: %d, Минуты: %d, Секунды: %d\n")
		// --------------------------------------------
		default:
			msg.Text = consts.TEXT_DEFAULt
		}

		// send  message to user
		if _, err := b.api.Send(msg); err != nil {
			b.logger.Error("faild to send message", "error", err)
		}

	}

}
