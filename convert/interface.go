package convert

import (
	"encoding/json"
	"fmt"
)

func InterfaceToBytes(v interface{}) ([]byte, error) {
	var vals []byte
	switch val := v.(type) {
	case string:
		vals = []byte(val)
	case []byte:
		vals = val
	case int64, int32, int, float64, float32:
		vals = []byte(fmt.Sprintf("%v", val))
	default:
		var err error
		vals, err = json.Marshal(val)
		if err != nil {
			return nil, err
		}
	}

	return vals, nil
}

func InterfaceToString(v interface{}) (string, error) {
	vals, err := InterfaceToBytes(v)
	return string(vals), err
}

func InterfaceToHashMapInterface(v interface{}) (map[string]interface{}, error) {
	vals, err := InterfaceToBytes(v)
	if err != nil {
		return nil, err
	}

	hashMaps := make(map[string]interface{})
	err = json.Unmarshal(vals, &hashMaps)
	if err != nil {
		return nil, err
	}

	return hashMaps, nil
}
