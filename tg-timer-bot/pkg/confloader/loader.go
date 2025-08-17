package confloader

import (
	"errors"
	"log/slog"
	"os"

	"github.com/Talos-hub/tg-bot-cooking-timer-/pkg/consts"
)

const NAME_ENV = "TOKEN"

// LoadConfig is function that we should use to get configuration
func LoadConfig(logger *slog.Logger) (*Config, error) {
	botToken := os.Getenv(NAME_ENV)

	if len(botToken) == 0 {
		logger.Error("required env variables are messing")
		return nil, errors.New("required env variables are messing")
	}

	i, err := LoadData(consts.DEFAULT_MEAT_PATH, consts.DEFAULT_EGG_PATH)
	if err != nil {
		logger.Error("couldn't getting data", "error", err)
		i = defaultInterval()
		UpdateOrCreateConfig(consts.DEFAULT_EGG_PATH, &i.Egg)
		UpdateOrCreateConfig(consts.DEFAULT_MEAT_PATH, &i.Meat)
	}

	if i == nil {
		return nil, errors.New("error, data is empty")
	}

	config := NewConfig(botToken, i)
	return config, nil

}
