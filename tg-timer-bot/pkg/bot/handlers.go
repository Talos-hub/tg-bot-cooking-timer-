package bot

import (
	"strconv"
	"strings"

	conf "github.com/Talos-hub/tg-bot-cooking-timer-/pkg/confloader"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handleUserInput is reaquired hander that we use that user could be set up a config
func (b *bot) handleUserInput(msg *tgbotapi.Message, state *UserState) {
	parts := strings.Fields(msg.Text) // split msg without spaces
	// check that values are correct
	if len(parts) != 3 {
		replay := tgbotapi.NewMessage(msg.Chat.ID, "Неверный формат, введите три числа: час, минута, секунда")

		if _, err := b.api.Send(replay); err != nil {
			b.logger.Error("faild to send message", "error", err)
		}
		return
	}

	// converting>
	hour, err1 := strconv.Atoi(parts[0])
	minute, err2 := strconv.Atoi(parts[1])
	second, err3 := strconv.Atoi(parts[2])
	// <convering

	// check err >
	if err1 != nil || err2 != nil || err3 != nil {
		replay := tgbotapi.NewMessage(msg.Chat.ID, "Все значения должны быть числами!!!")

		if _, err := b.api.Send(replay); err != nil {
			b.logger.Error("faild to send message", "error", err)
		}
		return
	}
	// chech < err

	// create interval then put it to IntervalFoodTime
	interval := conf.IntervalTime{
		Hours:  hour,
		Minute: minute,
		Second: second,
	}

	config := &conf.IntervalFoodTime{}

	// check food type
	if state.FoodType == "meat" {
		config.Meat = interval
	} else {
		config.Egg = interval
	}

	// convert int64 to string
	chat := strconv.FormatInt(msg.Chat.ID, 10)

	// create new path for config
	path := chat + conf.JSON_NAME
	err := conf.UpdateOrCreateConfig(path, *config) // creating or update config

	// check err >
	if err != nil {
		replay := tgbotapi.NewMessage(msg.Chat.ID, "Ошибка!!! Извините за неудобства, в ближайшее время будет исправленно.")
		b.logger.Error("faild creating or saving config", "error", err)

		if _, err := b.api.Send(replay); err != nil {
			b.logger.Error("faild to send message", "error", err)
		}
	}
	// check < err
}
