package env

import "strconv"

// GetInt64 returns integer64 from env variable
func GetInt64(key string, defaultValues ...int64) int64 {
	val, ok := getEnv(key)
	if !ok {
		if len(defaultValues) > 0 {
			return defaultValues[0]
		}
		
		return 0
	}
	
	valInt, err := strconv.ParseInt(val, 0, 64)
	if err != nil {
		if len(defaultValues) > 0 {
			return defaultValues[0]
		}
		
		return 0
	}
	
	return valInt
}
