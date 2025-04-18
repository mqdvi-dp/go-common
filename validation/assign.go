package validation

import (
	"fmt"
	"reflect"
	"strconv"
)

// assignValue set value into struct field
func (v *validation) assignValue(val string, rsf reflect.StructField, rv reflect.Value) error {
	kind := rsf.Type.Kind()

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.Atoi(val)
		if err != nil {
			return err
		}

		rv.SetInt(int64(i))
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}

		rv.SetFloat(f)
	case reflect.String:
		rv.SetString(val)
	case reflect.Bool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			b = false
		}

		rv.SetBool(b)
	default:
		return fmt.Errorf("data types %s cannot to binding", kind.String())
	}

	return nil
}
