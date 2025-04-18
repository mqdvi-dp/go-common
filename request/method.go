package request

import (
	"context"
	"net/http"
)

type MethodInterface interface {
	// Get is request with method GET
	Get(ctx context.Context) ([]byte, int, string, error)
	// Post is a request with method POST
	Post(ctx context.Context, payload []byte) ([]byte, int, string, error)
	// Put is request with method PUT
	Put(ctx context.Context, payload []byte) ([]byte, int, string, error)
	// Delete is request with method DELETE
	Delete(ctx context.Context, payload []byte) ([]byte, int, string, error)
}

func (r *request) Get(ctx context.Context) ([]byte, int, string, error) {
	return r.wrapper(ctx, nil, http.MethodGet)
}

func (r *request) Post(ctx context.Context, payload []byte) ([]byte, int, string, error) {
	return r.wrapper(ctx, payload, http.MethodPost)
}

func (r *request) Put(ctx context.Context, payload []byte) ([]byte, int, string, error) {
	return r.wrapper(ctx, payload, http.MethodPut)
}

func (r *request) Delete(ctx context.Context, payload []byte) ([]byte, int, string, error) {
	return r.wrapper(ctx, payload, http.MethodDelete)
}
