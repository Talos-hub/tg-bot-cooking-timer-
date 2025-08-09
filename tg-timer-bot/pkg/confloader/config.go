package confloader

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	// token telegram api
	Token string
	Food  IntervalFoodTime
}

type IntervalTime struct {
	Second int `json:"second"`
	Minute int `json:"minute"`
	Hours  int `json:"hours"`
}

type IntervalFoodTime struct {
	Meat IntervalTime `json:"meat"`
	Egg  IntervalTime `json:"egg"`
}

const jsonName = ".json"

// UpdateTimeConfig update or create new time configuration.
func UpdateTimeConfig(i IntervalFoodTime, path string) error {
	file, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error open or creating the file <<%s>>, err: %w", path, err)
	}
	defer file.Close()

	encdoer := json.NewEncoder(file)
	// for readadbility
	encdoer.SetIndent("", " ")

	err = encdoer.Encode(i)
	if err != nil {
		return fmt.Errorf("error endcoding intervaltime: %w", err)
	}

	return nil
}

// TODO
func DefaultInterval(c *Config) error {
	c.Food.Meat.Minute = 20
	c.Food.Egg.Minute = 8

	if c.Food.Egg.Minute != 8 {
		return errors.New("error: egg.Minute isn't correct")
	} else if c.Food.Meat.Minute != 20 {
		return errors.New("error: meat.Minute isn't correct")
	}

	return nil

}

// UserSetUp set up custom config
func UserSetUp(path string, i IntervalFoodTime) error {
	// validation path
	if len(path) == 0 {
		return fmt.Errorf("error, the path: <<%s>> less than zero", path)
	}

	js := jsonName
	lenPath := len(path)
	lenJs := len(js)

	// validation path
	if b := strings.Contains(path[:lenPath-lenJs], js); !b {
		return errors.New("error, this isn't json file")
	}

	// open file
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating the file <<%s>>, err: %w", path, err)
	}

	en := json.NewEncoder(file)
	// for readability
	en.SetIndent("", " ")

	err = en.Encode(i)
	if err != nil {
		return fmt.Errorf("error encoding: %w", err)
	}

	return nil
}

// IsExsitYserCustomConfig check if config is exist.
// It return true if file is exist and false if not
func IsExisttUserCustomConfig(chatId int) bool {
	strChatId := strconv.Itoa(chatId)
	path := fmt.Sprint(strChatId, jsonName)

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}
