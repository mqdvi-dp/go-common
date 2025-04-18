package request

import (
	"fmt"
	"reflect"
)

func validatePtr(val interface{}) error {
	vof := reflect.ValueOf(val)

	if vof.Kind() != reflect.Ptr {
		return fmt.Errorf("value should be pointer")
	}

	elem := vof.Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("value should be struct with pointer")
	}

	return nil
}
