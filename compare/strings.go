package compare

import "strings"

func String(val string, values ...string) bool {
	for _, v := range values {
		if strings.EqualFold(val, v) {
			return true
		}
	}

	return false
}
