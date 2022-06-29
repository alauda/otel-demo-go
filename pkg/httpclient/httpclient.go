package httpclient

import (
	"context"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
)

type HttpClient struct {
	client http.Client
}

func DefaultClient() HttpClient {
	return HttpClient{
		client: http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)},
	}
}

func (c HttpClient) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req)
}
