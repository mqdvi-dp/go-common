package version

import (
	"context"
	"fmt"

	"github.com/mqdvi-dp/go-common/tracer"
)

// LT method for checking appVersion has criteria for the maximal version or not
func (v *Version) LT(ctx context.Context, maxVersion string) bool {
	if v == nil {
		tracer.SetError(ctx, fmt.Errorf("app version not found"))
		return false
	}

	if len(v.AppVersion) < 1 {
		tracer.SetError(ctx, fmt.Errorf("app version not found"))
		return false
	}

	mv, err := convert(maxVersion)
	if err != nil {
		tracer.SetError(ctx, fmt.Errorf("error convert maxVersion %v", err))
		return false
	}

	val := compare(v.AppVersion, mv)
	return val < 0
}

// LTE method for checking appVersion has criteria for the maximal version or not
func (v *Version) LTE(ctx context.Context, maxVersion string) bool {
	if v == nil {
		tracer.SetError(ctx, fmt.Errorf("app version not found"))
		return false
	}

	if len(v.AppVersion) < 1 {
		tracer.SetError(ctx, fmt.Errorf("app version not found"))
		return false
	}

	mv, err := convert(maxVersion)
	if err != nil {
		tracer.SetError(ctx, fmt.Errorf("error convert maxVersion %v", err))
		return false
	}

	val := compare(v.AppVersion, mv)
	return val <= 0
}
