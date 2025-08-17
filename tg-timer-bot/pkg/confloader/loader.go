package confloader

import (
	"errors"
	"log/slog"
	"os"
)

const NAME_ENV = "TOKEN"

func LoadConfig(logger *slog.Logger) (*Config, error) {
	botToken := os.Getenv(NAME_ENV)

	if len(botToken) == 0 {
		logger.Error("required env variables are messing")
		return nil, errors.New("required env variables are messing")
	}

	i, err := LoadData(DEFAULT_MEAT_PATH, DEFAULT_EGG_PATH)
	if err != nil {
		logger.Error("error, couldn't getting data")
		i = defaultInterval()
	}

	if i != nil {
		return nil, errors.New("error, data is empty")
	}

	config := NewConfig(botToken, i)
	return config, nil

}
