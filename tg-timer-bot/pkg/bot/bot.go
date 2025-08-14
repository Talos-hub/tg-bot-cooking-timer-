package bot

import (
	"errors"
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// UserState is necessary structure that we use to
// check that a user is submitting data.
type UserState struct {
	WaitingForInput bool
	FoodType        string
}

// Bot is struct we use to work with telegram api
type Bot struct {
	api       *tgbotapi.BotAPI     // api
	userState map[int64]*UserState // check that a user is submitting data.
	logger    *slog.Logger         // logging errors
}

// NewBot is constructor that return a Bot struct
func NewBot(token string, logger *slog.Logger) (*Bot, error) {
	// check token length
	l := len(token)
	if l == 0 || l < 5 {
		return nil, fmt.Errorf("error, the token is not correct, length: %d", l)
	}
	// check that logger is not nil
	if logger == nil {
		return nil, errors.New("error no logger")
	}

	// create new bot api
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("error creating a telegram bot api: %w", err)
	}

	return &Bot{
		api:       bot,
		userState: map[int64]*UserState{},
		logger:    logger,
	}, nil

}
