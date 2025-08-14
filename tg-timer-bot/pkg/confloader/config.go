package confloader

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	JSON_NAME    = ".json"
	DEFAULT_PATH = "default.json" // this is path for default settings
)

// Config is configuration struct.
// You should use it when you're going to create telegram bot
type Config struct {
	// token telegram api
	Token string
	Food  IntervalFoodTime
}

// IntervalTime is struct that we need to set up time for our bot
type IntervalTime struct {
	Second int `json:"second"`
	Minute int `json:"minute"`
	Hours  int `json:"hours"`
}

// IntervalFoodTime is struct that include different type food
type IntervalFoodTime struct {
	Meat IntervalTime `json:"meat"`
	Egg  IntervalTime `json:"egg"`
}

// defultInterval is function that return a default settings for config.
// you should use it when you can't load or create a config.
func defaultInterval() *IntervalFoodTime {
	//meat
	m := IntervalTime{
		Second: 0,
		Hours:  0,
		Minute: 20,
	}
	// egg
	e := IntervalTime{
		Second: 0,
		Hours:  0,
		Minute: 8,
	}

	i := IntervalFoodTime{
		Meat: m,
		Egg:  e,
	}

	return &i
}

// UserSetUp set up a custom user config
func UpdateOrCreateConfig(path string, i IntervalFoodTime) error {
	// validation path
	if len(path) == 0 {
		return fmt.Errorf("error, the path: <<%s>> less than zero", path)
	}

	// for readability
	js := JSON_NAME

	lenPath := len(path)
	lenJs := len(js)

	// validation path
	// we checking the [:lenpath-lenJs] because we cut off the path
	// and checking the json extention
	if b := strings.Contains(path[lenPath-lenJs:], js); !b {
		return errors.New("error, this isn't json file")
	}

	// open file
	file, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error open or creating the file <<%s>>, err: %w", path, err)
	}
	defer file.Close()

	en := json.NewEncoder(file)
	// for readability
	en.SetIndent("", " ")

	err = en.Encode(i)
	if err != nil {
		return fmt.Errorf("error encoding: %w", err)
	}

	return nil
}

// IsExsitUserCustomConfig check if config is exist.
// It return true if a file is exist and false if not
func IsExisttUserCustomConfig(chatId int) bool {
	strChatId := strconv.Itoa(chatId)
	path := fmt.Sprint(strChatId, JSON_NAME)
	return checkPath(path)
}

// checkPath is simple function that check a file path
func checkPath(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// LoadIntreval read a ItervalFoodTime and return it.
// It's necessery in order to make a Config struct
func loadInterval() (*IntervalFoodTime, error) {
	// check that the default config is exist if not
	// then we create new config
	if p := checkPath(DEFAULT_PATH); !p {
		i := defaultInterval()
		UpdateOrCreateConfig(DEFAULT_PATH, *i)
		return i, nil
	}

	// for readability
	path := DEFAULT_PATH
	var i IntervalFoodTime

	// open a file
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("the file <<%s>>, is not exist, err: %w", path, err)
		}
		if os.IsPermission(err) {
			return nil, fmt.Errorf("error, permission denied, err: %w", err)
		}
		return nil, fmt.Errorf("error open the file <<%s>>, err: %w", path, err)
	}

	// new decoder
	d := json.NewDecoder(file)

	err = d.Decode(&i)
	if err != nil {
		return nil, fmt.Errorf("error decoding: %w", err)
	}

	return &i, nil
}

// NewConfig is constructor
func NewConfig(token string, i *IntervalFoodTime) *Config {
	return &Config{Token: token, Food: *i}
}
