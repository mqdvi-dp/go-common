package version

import (
	"errors"
	"fmt"
)

var (
	ErrVersionNotFound = errors.New("version not found")
)

type Version struct {
	AppVersion []int
}

// New is construct will return an instance for checking Application Version
func New(appVersion string) (*Version, error) {
	v, err := convert(appVersion)
	if err != nil {
		return nil, err
	}

	return &Version{
		AppVersion: v,
	}, nil
}

func (v *Version) String() string {
	if v == nil {
		return ""
	}

	return fmt.Sprintf("%d.%d.%d", v.AppVersion[0], v.AppVersion[1], v.AppVersion[2])
}
