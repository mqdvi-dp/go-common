package version

import (
	"context"
	"fmt"

	"github.com/mqdvi-dp/go-common/tracer"
)

// GT method for checking appVersion has criteria for the minimal version or not
func (v *Version) GT(ctx context.Context, minVersion string) bool {
	if v == nil {
		tracer.SetError(ctx, fmt.Errorf("app version not found"))
		return false
	}

	if len(v.AppVersion) < 1 {
		tracer.SetError(ctx, fmt.Errorf("app version not found"))
		return false
	}

	mv, err := convert(minVersion)
	if err != nil {
		tracer.SetError(ctx, fmt.Errorf("error convert minVersion %v", err))
		return false
	}

	val := compare(v.AppVersion, mv)
	return val > 0
}

// GTE method for checking appVersion has criteria for the minimal version or not
func (v *Version) GTE(ctx context.Context, minVersion string) bool {
	if v == nil {
		tracer.SetError(ctx, fmt.Errorf("app version not found"))
		return false
	}

	if len(v.AppVersion) < 1 {
		tracer.SetError(ctx, fmt.Errorf("app version not found"))
		return false
	}

	mv, err := convert(minVersion)
	if err != nil {
		tracer.SetError(ctx, fmt.Errorf("error convert minVersion %v", err))
		return false
	}

	val := compare(v.AppVersion, mv)
	return val >= 0
}
