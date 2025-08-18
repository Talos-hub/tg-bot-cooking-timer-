package bot

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	conf "github.com/Talos-hub/tg-bot-cooking-timer-/pkg/confloader"
	"github.com/Talos-hub/tg-bot-cooking-timer-/pkg/consts"
	"github.com/Talos-hub/tg-bot-cooking-timer-/pkg/paths"
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

	// creating or update a config
	path, err := paths.CreateNewPath(msg.Chat.ID, state.FoodType)
	if err != nil {
		b.logger.Error("error handling user input", "error", err)
	}
	err = conf.UpdateOrCreateConfig(path, &interval)

	// check err >
	if err != nil {
		replay := tgbotapi.NewMessage(msg.Chat.ID, "Ошибка!!! Извините за неудобства, в ближайшее время будет исправленно.")
		b.logger.Error("faild creating or saving config", "error", err)
		if _, err := b.api.Send(replay); err != nil {
			b.logger.Error("faild to send message", "error", err)
		}
	}
	// check < err

	replay := tgbotapi.NewMessage(msg.Chat.ID, "Найстройки сохранены успешно")
	if _, err := b.api.Send(replay); err != nil {
		b.logger.Error("faild to send message", "error", err)
	}
}

// ShowSettings retruns a string with data about current settings
func (b *bot) ShowSettings(chatID int64, defaultSettings *conf.IntervalFoodTime) (string, error) {
	ok := conf.IsExistUserConfigs(chatID)
	msg := ""
	// default settings
	if !ok {
		// update msg and then returns it
		msg = fmt.Sprintf(
			"Настройки для Мяса: Часы: %d, Минуты: %d, Секунды: %d\nНастройки для Яйца: Часы: %d, Минуты: %d, Секунды: %d\n",
			defaultSettings.Meat.Hours, defaultSettings.Meat.Minute, defaultSettings.Meat.Second,
			defaultSettings.Egg.Hours, defaultSettings.Egg.Minute, defaultSettings.Egg.Second,
		)
		return msg, nil // after don't forget remove it
	} else { // custom settings
		meat, err1 := paths.CreateNewPath(chatID, "meat")
		egg, err2 := paths.CreateNewPath(chatID, "egg")

		if err1 != nil || err2 != nil {
			return "", fmt.Errorf("error showing, %w, %w", err1, err2)
		}

		data, err := conf.LoadData(meat, egg)
		if err != nil {
			return "", fmt.Errorf("error showing, %w", err)
		}
		// update msg and then retruns it
		msg = fmt.Sprintf(
			"Настройки для Мяса: Часы: %d, Минуты: %d, Секунды: %d\nНастройки для Яйца: Часы: %d, Минуты: %d, Секунды: %d\n",
			data.Meat.Hours, data.Meat.Minute, data.Meat.Second,
			data.Egg.Hours, data.Egg.Minute, data.Egg.Second,
		)
	}

	// retruns values
	return msg, nil
}

// StarTimer is fucntion that implement logic of a timer
func (b *bot) StartTimer(chatID int64, typeFood string, config conf.IntervalFoodTime) (string, error) {
	ok := conf.IsExistUserConfigs(chatID)

	// If custom confgis are exist then defaultConfgis = CastomConfigs
	if ok {
		meat, err := paths.CreateNewPath(chatID, consts.MEAT)
		if err != nil {
			return "", fmt.Errorf("error starting: %w", err)
		}
		egg, err := paths.CreateNewPath(chatID, consts.EGG)
		if err != nil {
			return "", fmt.Errorf("error starting: %w", err)
		}
		data, err := conf.LoadData(meat, egg)
		if err != nil {
			return "", err
		}
		// Now Default = Castom
		config = *data
	}

	var duration time.Duration

	// check food type
	switch typeFood {
	case consts.MEAT: // meat
		duration = time.Duration(config.Meat.Hours)*time.Hour + // creating duration
			time.Duration(config.Meat.Minute)*time.Minute + time.Duration(config.Meat.Second)*time.Second
	case consts.EGG: // egg
		duration = time.Duration(config.Egg.Hours)*time.Hour + // creating duration
			time.Duration(config.Egg.Minute)*time.Minute + time.Duration(config.Egg.Second)*time.Second
	}

	go func() {
		time.Sleep(duration)
		msg := tgbotapi.NewMessage(chatID, "Время приготовления Истекло!!!")

		for i := 0; i < 3; i++ {
			if _, err := b.api.Send(msg); err != nil {
				b.logger.Error("faild to send message", "error", err)
			}
		}

	}()

	return "Таймер зампушен на " + duration.String(), nil
}
