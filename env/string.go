package env

// GetString returns string from env variable
func GetString(key string, defaultValues ...string) string {
	val, ok := getEnv(key)
	if !ok {
		// when default values are exists
		if len(defaultValues) > 0 {
			return defaultValues[0]
		}
		
		return ""
	}
	
	return val
}
