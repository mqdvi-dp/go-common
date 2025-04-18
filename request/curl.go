package request

import (
	"bytes"
	"io"
	"net/http"

	"moul.io/http2curl"
)

func (r *request) curl(payload []byte, method string) ([]byte, http.Header, int, string, error) {
	req, err := http.NewRequest(method, r.url, buf(payload))
	if err != nil {
		return nil, nil, http.StatusInternalServerError, "", err
	}

	// create header
	if r.header != nil {
		req.Header = r.header
	}

	// set basic auth if exists
	if r.basicAuth != nil && r.basicAuth.set {
		req.SetBasicAuth(r.basicAuth.username, r.basicAuth.password)
	}

	// Get curl command
	command, err := http2curl.GetCurlCommand(req)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, "", err
	}

	// do request to a client
	res, err := r.client.Do(req)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, command.String(), err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, command.String(), err
	}

	return body, res.Header, res.StatusCode, command.String(), nil
}

func buf(p []byte) io.ReadCloser {
	if p != nil {
		r := bytes.NewReader(p)
		return io.NopCloser(r)
	}

	return nil
}
