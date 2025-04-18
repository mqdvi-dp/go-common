package validation

import (
	"context"

	"github.com/mqdvi-dp/go-common/tracer"
)

func (v *validation) Validate(ctx context.Context, dest interface{}) error {
	err := v.validate.Struct(dest)
	if err != nil {
		tracer.SetError(ctx, err)

		return err
	}

	return nil
}
