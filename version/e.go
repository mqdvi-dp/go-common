package version

import (
	"context"
	"fmt"

	"github.com/mqdvi-dp/go-common/tracer"
)

// Eq method for checking appVersion has criteria for an equal version or not
func (v *Version) Eq(ctx context.Context, version string) bool {
	if v == nil {
		tracer.SetError(ctx, fmt.Errorf("app version not found"))
		return false
	}

	if len(v.AppVersion) < 1 {
		tracer.SetError(ctx, fmt.Errorf("app version not found"))
		return false
	}

	mv, err := convert(version)
	if err != nil {
		tracer.SetError(ctx, fmt.Errorf("error convert version %v", err))
		return false
	}

	val := compare(v.AppVersion, mv)
	return val == 0
}
