package env

import "strconv"

// GetFloat64 returns float64 (double) from env variable
func GetFloat64(key string, defaultValues ...float64) float64 {
	val, ok := getEnv(key)
	if !ok {
		if len(defaultValues) > 0 {
			return defaultValues[0]
		}
		
		return 0
	}
	
	valFloat, err := strconv.ParseFloat(val, 64)
	if err != nil {
		if len(defaultValues) > 0 {
			return defaultValues[0]
		}
		
		return 0
	}
	
	return valFloat
}
