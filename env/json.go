package env

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func GetJSON(key string, destination interface{}, isRecrusive ...bool) error {
	// get value from interface
	vof := reflect.ValueOf(destination)
	// check value interface, should be pointer
	if vof.Kind() != reflect.Ptr && len(isRecrusive) > 0 && !isRecrusive[0] {
		return fmt.Errorf("destination should be pointer")
	}

	val, ok := getEnv(key)
	if !ok {
		return fmt.Errorf("key '%s' not found in env variable", key)
	}

	err := json.Unmarshal([]byte(val), destination)
	if err != nil {
		return err
	}

	return nil
}
