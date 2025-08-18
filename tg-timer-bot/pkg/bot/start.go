package bot

import (
	conf "github.com/Talos-hub/tg-bot-cooking-timer-/pkg/confloader"
	"github.com/Talos-hub/tg-bot-cooking-timer-/pkg/consts"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Start is launch the bot
func (b *bot) Start(defaultConf *conf.IntervalFoodTime) {
	if defaultConf == nil {
		b.logger.Error("there is no default config")
		return
	}

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
				FoodType:        consts.MEAT,
			}
		// --------------------------------------
		case consts.EGG: // setup settings for egg
			msg.Text = consts.TEXT_EGG
			// change user state
			b.userState[chatID] = &UserState{
				WaitingForInput: true,
				FoodType:        consts.EGG,
			}
		// --------------------------------------------
		case consts.START_TIMER: // start timer
			msg.Text = consts.TEXT_START_TIMER
		// --------------------------------------------
		case consts.SHOW:
			m, err := b.ShowSettings(chatID, defaultConf)
			if err != nil {
				b.logger.Error("error showing settings", "error", err)
				msg.Text = "Извените за неудобства, функция не временно не работает"
			} else {
				msg.Text = m
			}
		// --------------------------------------------
		case consts.EGG_TIMER: // egg timer
			m, err := b.StartTimer(chatID, consts.EGG, *defaultConf)
			if err != nil {
				b.logger.Error("error starting timer", "error", err)
			}
			msg.Text = m
		// --------------------------------------------
		case consts.MEAT_TIMER: // meat timer
			m, err := b.StartTimer(chatID, consts.MEAT, *defaultConf)
			if err != nil {
				b.logger.Error("error startin timer", "error", err)
			}
			msg.Text = m
		default:
			msg.Text = consts.TEXT_DEFAULt
		}

		// send  message to user
		if _, err := b.api.Send(msg); err != nil {
			b.logger.Error("faild to send message", "error", err)
		}

	}

}
