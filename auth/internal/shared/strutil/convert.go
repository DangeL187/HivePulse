package strutil

import (
	"errors"
	"strconv"

	"github.com/DangeL187/erax"
)

func StringToUint(s string) (uint, error) {
	if s == "" {
		return 0, errors.New("empty string cannot be converted to uint")
	}

	parsed, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, erax.Wrap(err, "failed parse uint")
	}

	return uint(parsed), nil
}
