package env

import "strconv"

// GetInt returns int from env variable
func GetInt(key string, defaultValues ...int) int {
	val, ok := getEnv(key)
	if !ok {
		if len(defaultValues) > 0 {
			return defaultValues[0]
		}
		
		return 0
	}
	
	valInt, err := strconv.Atoi(val)
	if err != nil {
		if len(defaultValues) > 0 {
			return defaultValues[0]
		}
		
		return 0
	}
	
	return valInt
}
