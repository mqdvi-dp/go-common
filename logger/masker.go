package logger

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/ggwhite/go-masker"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/helpers"
)

var mapMaskTypes = map[string]func(string) string{
	"password":      maskPassword,
	"pin":           maskPassword,
	"otp":           maskPassword,
	"phone_number":  maskPhoneNumber,
	"email":         maskEmail,
	"token":         maskToken,
	"refreshtoken":  maskToken,
	"authorization": maskToken,
}

func splitAndMap(str string, separator string, mapFunc func(string) string) string {
	parts := strings.Split(str, separator)

	for i, part := range parts {
		parts[i] = mapFunc(part)
	}

	return strings.Join(parts, separator)
}

func maskPassword(val string) string {
	return masker.Password(val)
}

func maskPhoneNumber(val string) string {
	return masker.Mobile(val)
}

func maskEmail(val string) string {
	return masker.Email(val)
}

func maskValue(val string) string {
	return splitAndMap(val, " ", masker.Name)
}

func maskToken(val string) string {
	return splitAndMap(val, " ", helpers.SHA256)
}

func MaskedCredentials(b []byte) []byte {
	maskers := env.GetListString("DATA_MASKED", "password", "pin", "email", "phone_number", "username", "token", "authorization", "otp")

	// dd means dynamic-data
	var dd = make(map[string]interface{})
	err := json.Unmarshal(b, &dd)
	if err != nil {
		return b
	}

	for _, mask := range maskers {
		value, ok := dd[mask]
		if !ok {
			continue
		}

		// avoid panic error when value is nil
		if value == nil {
			continue
		}
		vof := reflect.ValueOf(value)

		var val string
		switch vof.Type().Kind() {
		case reflect.Float64, reflect.Float32:
			val = strconv.FormatFloat(vof.Float(), 'f', 0, 64)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val = strconv.Itoa(int(vof.Int()))
		case reflect.String:
			val = vof.String()
		case reflect.Slice:
			array := make([]interface{}, 0)
			array = append(array, vof.Interface())
			val = fmt.Sprintf("%v", array)
		default:
			continue
		}

		// because username can be as email and phone_number
		// we need to check the value first
		if strings.EqualFold(mask, "username") {
			if strings.Contains(val, "@") {
				// will mask as email
				dd[mask] = maskEmail(val)
				continue
			}

			// default is phone_number
			dd[mask] = maskPhoneNumber(val)
			continue
		}

		mf, ok := mapMaskTypes[mask]
		if !ok {
			// default
			dd[mask] = maskValue(val)
			continue
		}

		dd[mask] = mf(val)
	}

	// marshal the data again
	by, err := json.Marshal(dd)
	if err != nil {
		return b
	}

	return by
}
