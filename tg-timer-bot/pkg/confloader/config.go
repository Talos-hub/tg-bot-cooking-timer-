package confloader

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/Talos-hub/tg-bot-cooking-timer-/pkg/consts"
	"github.com/Talos-hub/tg-bot-cooking-timer-/pkg/paths"
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
func UpdateOrCreateConfig(path string, i *IntervalTime) error {
	// validation path
	if err := paths.ValidationJsonPath(path); err != nil {
		return err
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
func CheckUserConfigFile(chatId int, typeFood string) bool {
	path, err := paths.CreateNewPath(int64(chatId), typeFood)
	if err != nil {
		return false
	}
	return checkPath(path)
}

// IsExsitUserCustomConfig check if two configs are exist.
// It return true if a file are exist and false if not
func IsExistUserConfigs(chatID int64) bool {
	wg := sync.WaitGroup{}

	ok1 := make(chan bool)
	ok2 := make(chan bool)

	// after works close channels
	defer func() {
		close(ok1)
		close(ok2)
	}()

	f := func(wg *sync.WaitGroup, ok chan bool, typeFood string) {
		ok <- CheckUserConfigFile(int(chatID), typeFood) // start searching path
		wg.Done()
	}
	wg.Add(2)
	go f(&wg, ok1, consts.MEAT)
	go f(&wg, ok2, consts.EGG)

	meat := <-ok1
	egg := <-ok2
	wg.Wait()

	if !meat || !egg {
		return false
	}
	return true

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

// LoadData load json data and return pointer on IntervalFoodTime.
// first and second path should be correct format either like name + .json or like ChatID + TypeFood + .json
func LoadData(meatPath, eggPath string) (*IntervalFoodTime, error) {
	err1 := paths.ValidationJsonPath(meatPath)
	err2 := paths.ValidationJsonPath(eggPath)

	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("error these files or paths aren't correct: %w, %w", err1, err2)
	}
	// open files
	file1, err1 := os.Open(meatPath)
	file2, err2 := os.Open(eggPath)

	// check errors
	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("error open these files: %w, %w", err1, err2)
	}

	//======================difficult to read=========================

	w := sync.WaitGroup{}
	// for first file
	timeChan := make(chan IntervalTime)
	errChan := make(chan error)
	// for second file
	timeChan2 := make(chan IntervalTime)
	errChan2 := make(chan error)

	// start

	w.Add(2)
	go loadIntervalTime(file1, errChan, timeChan, &w)
	go loadIntervalTime(file2, errChan2, timeChan2, &w)

	// gets errors
	err1 = <-errChan
	err2 = <-errChan2
	// gets interval
	meat := <-timeChan
	egg := <-timeChan2

	// wait while all goroutines is done
	w.Wait()

	//close channels
	defer func() {
		close(errChan)
		close(errChan2)
		close(timeChan)
		close(timeChan2)
	}()

	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("error read data: %w, %w", err1, err2)
	}

	//==========================================================================

	return &IntervalFoodTime{
		Meat: meat,
		Egg:  egg,
	}, nil

}

// loadIntervalTime is function that read a json data form a file and send it to timeChan,
// if it can't read a data so it send an error to errChan and we have to handle the error.
// Also the function close a file when read is done
func loadIntervalTime(r *os.File, errChan chan error, timeChan chan IntervalTime, w *sync.WaitGroup) {
	defer r.Close()
	decoder := json.NewDecoder(r)

	var interval IntervalTime

	// if a err is exist so we send the error to a errChan and up level we have to handle it
	err := decoder.Decode(&interval)
	if err != nil {
		errChan <- fmt.Errorf("error decoding: %w", err)
		return
	} else {
		errChan <- nil
	}
	// send data to a chan
	timeChan <- interval
	w.Done()
}

// NewConfig is constructor
func NewConfig(token string, i *IntervalFoodTime) *Config {
	return &Config{Token: token, Food: *i}
}
