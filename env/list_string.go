package env

import "strings"

// GetListString returns list string from env variable
func GetListString(key string, defaultValues ...string) []string {
	val, ok := getEnv(key)
	if !ok {
		// when default values are exists
		if len(defaultValues) > 0 {
			return defaultValues
		}
		
		return []string{}
	}
	
	return strings.Split(val, ",")
}
