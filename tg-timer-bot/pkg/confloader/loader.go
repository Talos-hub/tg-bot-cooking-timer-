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

	i, err := loadInterval()
	if err != nil {
		logger.Error("error, couldn't getting data")

		return nil, err
	}

	config := NewConfig(botToken, i)
	return config, nil

}
