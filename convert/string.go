package convert

import (
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
)

func StringWithUnderscored(str string) string {
	// replace all dashed with underscored
	str = strings.ReplaceAll(str, "-", "_")
	// replace all space with underscored
	str = strings.ReplaceAll(str, " ", "_")

	return strings.ToLower(str)
}

func StringToTitle(str string) string {
	// replace underscore with space
	str = strings.ReplaceAll(str, "_", " ")
	// replace any dot char with space
	str = strings.ReplaceAll(str, ".", " ")

	// force to lower
	str = strings.ToLower(str)

	return strings.ToTitle(str)
}

func FloatToString(val float64) string {
	return strconv.FormatFloat(val, 'f', 0, 64)
}

func IntToString(val int64) string {
	return strconv.FormatInt(val, 10)
}

func HumanizeNumber[T Number](val T, separators ...string) string {
	v := float64(val)

	// default separator is dot
	separator := "."
	if len(separators) > 0 {
		separator = separators[0]
	}

	return strings.ReplaceAll(humanize.Commaf(v), ",", separator)
}
