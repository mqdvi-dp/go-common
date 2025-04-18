package regex

import "regexp"

func ExtractNumber(s string) string {
	re := regexp.MustCompile(`\d+(\.\d+)*`)
	match := re.FindString(s)
	return match
}
