package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/mqdvi-dp/go-common/convert"
)

func reflection(v interface{}) (string, error) {
	vof := reflect.ValueOf(v)
	switch vof.Kind() {
	case reflect.String:
		return vof.String(), nil
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%v", vof.Interface()), nil
	case reflect.Ptr, reflect.Struct:
		return convert.InterfaceToString(vof.Interface())
	case reflect.Bool:
		return strconv.FormatBool(vof.Bool()), nil
	case reflect.Slice:
		var values []string
		for i := 0; i < vof.Len(); i++ {
			val, err := reflection(vof.Index(i).Interface())
			if err != nil {
				return "", err
			}

			values = append(values, val)
		}

		return strings.Join(values, ","), nil
	case reflect.Map:
		return convert.InterfaceToString(vof.Interface())
	default:
		return "", fmt.Errorf("kind not yet handled by reflection. Kind is %s", vof.Kind())
	}
}
