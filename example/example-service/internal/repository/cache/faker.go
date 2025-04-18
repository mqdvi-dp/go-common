package cache

import (
	"context"

	"github.com/mqdvi-dp/go-common/convert"
	"github.com/mqdvi-dp/go-common/example/example-service/internal/model"
	"github.com/mqdvi-dp/go-common/tracer"
)

func (c *cacheRepository) GetFaker(ctx context.Context) (resp model.ResponseFaker, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "CacheRepository:GetFaker")
	defer trace.Finish()

	key := prefixKeyFaker
	err = c.client.GetStruct(ctx, &resp, key)
	if err != nil {
		trace.SetError(err)
		return
	}

	return
}

func (c *cacheRepository) SetFaker(ctx context.Context, data model.ResponseFaker) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "CacheRepository:SetFaker")
	defer trace.Finish()

	key := prefixKeyFaker
	d, err := convert.InterfaceToString(data)
	if err != nil {
		trace.SetError(err)
		return
	}

	err = c.client.Set(ctx, key, d)
	if err != nil {
		trace.SetError(err)
		return
	}

	return
}
