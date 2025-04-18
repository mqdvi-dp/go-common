package env

import "strconv"

// GetBool returns boolean from env variable
func GetBool(key string, defaultValues ...bool) bool {
	val, ok := getEnv(key)
	if !ok {
		if len(defaultValues) > 0 {
			return defaultValues[0]
		}
		
		return false
	}
	
	valBoolean, err := strconv.ParseBool(val)
	if err != nil {
		if len(defaultValues) > 0 {
			return defaultValues[0]
		}
		
		return false
	}
	
	return valBoolean
}
