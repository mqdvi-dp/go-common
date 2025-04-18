package tracer

import (
	"encoding/json"
	"fmt"

	"github.com/mqdvi-dp/go-common/env"
)

// Option for init tracer options
type Option struct {
	ServiceName    string
	DSN            string
	Level          string
	BuildNumberTag string
	URLPath        string
	platformType   PlatformType
}

// defaultOption tracer is Jaeger
var defaultOption = Option{
	DSN:            env.GetString("JAEGER_HOST", "127.0.0.1:5775"),
	Level:          env.GetString("APP_LEVEL", "development"),
	BuildNumberTag: env.GetString("APP_VERSION", "v1.0.0"),
	platformType:   Jaeger,
}

// OptionFunc method
type OptionFunc func(*Option)

// SetDSN set agent host on options
func SetDSN(host string) OptionFunc {
	return func(option *Option) {
		option.DSN = host
	}
}

// SetLevel set level on options
func SetLevel(level string) OptionFunc {
	return func(option *Option) {
		option.Level = level
	}
}

// SetBuildNumberTag set build number tag on Option
func SetBuildNumberTag(number string) OptionFunc {
	return func(option *Option) {
		option.BuildNumberTag = number
	}
}

// SetPlatformType set platform type on Option
func SetPlatformType(pt PlatformType) OptionFunc {
	return func(option *Option) {
		option.platformType = pt
	}
}

// SetURLPath set the endpoint path of tracer platform
// this is actually for Jaeger, because otel library didn't support
// exporters Jaeger anymore since July 2023, so we need to use otelhttp
// instead of exporters/jaeger
func SetURLPath(urlPath string) OptionFunc {
	return func(o *Option) {
		o.URLPath = urlPath
	}
}

func toValue(v interface{}) string {
	var str string
	switch val := v.(type) {

	case uint, uint64, int, int64, float32, float64:
		str = fmt.Sprintf("%v", val)
	case error:
		if val != nil {
			str = val.Error()
		}
	case string:
		str = val
	case []byte:
		str = string(val)
	default:
		b, _ := json.Marshal(val)
		str = string(b)
	}

	return str
}
