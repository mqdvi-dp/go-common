package convert

import (
	"fmt"
	"strconv"
)

// StringToFloat convert string to float64
func StringToFloat(val string) float64 {
	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return -1
	}

	return v
}

// ToFloat convert to float from any types
func ToFloat(val interface{}) float64 {
	switch v := val.(type) {
	case int, int32, int64:
		return StringToFloat(fmt.Sprintf("%d", v))
	case string:
		return StringToFloat(v)
	default:
		return -1
	}
}
