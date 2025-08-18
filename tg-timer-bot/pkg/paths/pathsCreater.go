package paths

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Talos-hub/tg-bot-cooking-timer-/pkg/consts"
)

// CreateNewPath creat new path. If chatId == 0, it returns default path,
// if not it retruns = chatID + typeFood + .json
func CreateNewPath(chatID int64, typeFood string) (string, error) {
	return creatNewPath(chatID, typeFood)
}

// createNewPath creat new path. If chatId == 0, it returns default path,
// if not it retruns = chatID + typeFood + .json
func creatNewPath(chatID int64, typeFood string) (string, error) {
	if chatID == 0 {
		switch typeFood {
		case consts.EGG:
			newPath := consts.DEFAULT_EGG_PATH
			if err := ValidationJsonPath(newPath); err != nil {
				return "", err
			}
			return newPath, nil
		case consts.MEAT:
			newPath := consts.DEFAULT_MEAT_PATH
			if err := ValidationJsonPath(newPath); err != nil {
				return "", err
			}
			return newPath, nil
		default:
			return "", errors.New("unknown type food")

		}
	}

	str, err := ConvertInt64ToStr(chatID)
	if err != nil {
		return "", fmt.Errorf("couldn't convert value <<%d>>: %w", chatID, err)
	}

	lowerFood := strings.ToLower(typeFood)
	if i := strings.Compare(lowerFood, consts.MEAT); i == 0 {
		newPath := str + lowerFood + consts.JSON_NAME
		if err := ValidationJsonPath(newPath); err != nil {
			return "", err
		}
		return newPath, nil
	}

	if i := strings.Compare(lowerFood, consts.EGG); i == 0 {
		newPath := str + lowerFood + consts.JSON_NAME
		if err := ValidationJsonPath(newPath); err != nil {
			return "", err
		}
		return newPath, nil
	}

	return "", errors.New("unknown type food")
}

// ConferInt64ToStr is simple function that convert int64 to string
func ConvertInt64ToStr(num int64) (string, error) {
	str := strconv.FormatInt(num, 10)
	if len(str) == 0 {
		return "", errors.New("value isn't correct")
	}
	return str, nil
}
