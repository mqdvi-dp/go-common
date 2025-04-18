package env

import (
	"reflect"
	"time"
)

type TimeParam struct {
	Layout       string
	DefaultValue time.Time
}

// GetTime returns time.Time from env variable
func GetTime(key string, t TimeParam) time.Time {
	val, ok := getEnv(key)
	if !ok {
		return t.DefaultValue
	}
	
	// default value for layout time
	layout := time.RFC3339
	if !reflect.ValueOf(t.Layout).IsZero() {
		layout = t.Layout
	}
	
	vt, err := time.Parse(layout, val)
	if err != nil {
		return t.DefaultValue
	}
	
	return vt
}
