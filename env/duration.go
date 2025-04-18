package env

import "time"

// GetDuration returns time.Duration from env variable
func GetDuration(key string, defaultValues ...time.Duration) time.Duration {
	val, ok := getEnv(key)
	if !ok {
		// when default values is exists, return default values
		if len(defaultValues) > 0 {
			return defaultValues[0]
		}
		
		return time.Duration(0)
	}
	
	// parse to duration
	vd, err := time.ParseDuration(val)
	if err != nil {
		// when default values is exists, return default values
		if len(defaultValues) > 0 {
			return defaultValues[0]
		}
		
		return time.Duration(0)
	}
	
	return vd
}
