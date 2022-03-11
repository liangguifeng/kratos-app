package helper

import (
	"context"
	"io"
	"net/http"
	neturl "net/url"
	"time"
)

func NewHTTPClient() *httpClient {
	netClient := &http.Client{Timeout: time.Second * 90}
	return &httpClient{client: netClient}
}

type httpClient struct {
	client *http.Client
}

func (d *httpClient) GetClient() *http.Client {
	return d.client
}

func (d *httpClient) HttpGet(ctx context.Context, url string) (*http.Response, error) {
	resp, err := d.client.Get(url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (d *httpClient) HttpPost(ctx context.Context, url string, bodyType string, body io.Reader) (*http.Response, error) {
	resp, err := d.client.Post(url, bodyType, body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (d *httpClient) PostForm(ctx context.Context, url string, data neturl.Values) (*http.Response, error) {
	resp, err := d.client.PostForm(url, data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (d *httpClient) HttpHead(ctx context.Context, url string) (*http.Response, error) {
	resp, err := d.client.Head(url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (d *httpClient) HttpDo(ctx context.Context, r *http.Request) (*http.Response, error) {
	resp, err := d.client.Do(r)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
