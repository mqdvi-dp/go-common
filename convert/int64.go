package convert

import (
	"fmt"
	"strconv"
)

func StringToInt(val string) int64 {
	value, err := strconv.ParseInt(val, 0, 64)
	if err != nil {
		return -1
	}

	return value
}

func NumberToChar(i int) string {
	return fmt.Sprintf("%c", 'A'-1+i)
}
