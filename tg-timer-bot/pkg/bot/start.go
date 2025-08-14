package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

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
			// TODO, HANDLE that
			delete(b.userState, chatID)
			continue
		}

		// handle command
		switch update.Message.Command() {
		case START, HELP:
			msg.Text = ``
		}

	}

}
