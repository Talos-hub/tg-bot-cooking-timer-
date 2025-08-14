package bot

import (
	"errors"
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserState struct {
	WaitingForInput bool
	FoodType        string
}

type Bot struct {
	api       *tgbotapi.BotAPI
	userState map[int64]*UserState
	logger    *slog.Logger
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
