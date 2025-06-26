package not_hot

import (
	"errors"
	"regexp"
	"strconv"
)

var (
	idRegex = regexp.MustCompile(`"id":(\d+)`)
)

func parseDSId(str string) (int, error) {
	matches := idRegex.FindStringSubmatch(str)
	if len(matches) < 1 {
		return 0, errors.New("no id")
	}

	return strconv.Atoi(matches[1])
}
