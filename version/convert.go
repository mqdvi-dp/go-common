package version

import (
	"reflect"
	"strconv"
	"strings"
)

func convert(version string) ([]int, error) {
	if len(version) < 1 || reflect.ValueOf(version).IsZero() || version == "" {
		return nil, ErrVersionNotFound
	}
	var appVersion = make([]int, 3)
	s := strings.SplitN(version, ".", 3)

	for i, val := range s[:3] {
		v, err := strconv.Atoi(val)
		if err != nil {
			return nil, err
		}
		appVersion[i] = v
	}

	return appVersion, nil
}
