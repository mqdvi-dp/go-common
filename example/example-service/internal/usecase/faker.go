package usecase

import (
	"context"

	"github.com/mqdvi-dp/go-common/example/example-service/internal/model"
	"github.com/mqdvi-dp/go-common/tracer"
)

func (u *usecaseInstance) GetFaker(ctx context.Context) (resp model.ResponseFaker, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "Usecase:GetFaker")
	defer trace.Finish()

	// resp, err = u.cache.GetFaker(ctx)
	// if err != nil {
	resp, err = u.api.GetFaker(ctx)
	if err != nil {
		trace.SetError(err)
		return
	}

	// 	_ = u.cache.SetFaker(ctx, resp)
	// }

	return
}
