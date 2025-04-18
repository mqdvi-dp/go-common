package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/example/example-service/internal/model"
	"github.com/mqdvi-dp/go-common/tracer"
)

func (a *apiRepository) GetFaker(ctx context.Context) (resp model.ResponseFaker, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ApiRepository:GetFaker")
	defer trace.Finish()

	endpoint := fmt.Sprintf("%s%s", env.GetString("FAKER_HOST"), env.GetString("FAKER_PERSON"))
	header := http.Header{
		"Content-Type": []string{"application/json"},
	}

	aa := url.Values{}
	aa.Add("fcm", "asdbasdj")

	client := a.client.Request(header, "get_faker", endpoint)
	res, sc, _, err := client.Get(ctx)
	if err != nil {
		trace.SetError(err)

		return
	}

	if sc >= http.StatusBadRequest {
		err = fmt.Errorf("actual response %s", string(res))
		trace.SetError(err)

		return
	}

	err = json.Unmarshal(res, &resp)
	if err != nil {
		trace.SetError(err)

		return
	}

	return
}
