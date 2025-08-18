package paths

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Talos-hub/tg-bot-cooking-timer-/pkg/consts"
)

// ValudationJsonPath is validate that a path is correct and a file is json
func ValidationJsonPath(path string) error {

	if len(path) == 0 {
		return fmt.Errorf("error, the path: <<%s>> less than zero", path)
	}

	// for readability
	js := consts.JSON_NAME
	lenPath := len(path)
	lenJs := len(js)

	// validation path
	// we checking the [:lenpath-lenJs] because we cut off the path
	// and checking the json extention
	if b := strings.Contains(path[lenPath-lenJs:], js); !b {
		return errors.New("error, this isn't json file")
	}

	return nil
}
