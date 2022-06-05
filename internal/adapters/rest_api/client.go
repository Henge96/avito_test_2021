package rest_api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/shopspring/decimal"
)

type Client struct {
	client   http.Client
	apiKey   string
	basePath string
}

func New(apikey, basePath string) *Client {
	client := http.Client{
		Timeout: timeout,
	}

	return &Client{
		client:   client,
		apiKey:   apikey,
		basePath: basePath,
	}
}

func (c *Client) request(ctx context.Context, url, method string, body []byte) (decimal.Decimal, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, io.LimitReader(bytes.NewBuffer(body), 1048576))
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("c.client.Do: %w", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("io.ReadAll: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return decimal.Decimal{}, errUnknown
	}

	result := &Response{}
	err = json.Unmarshal(respBody, result)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("json.Unmarshall: %w", err)
	}

	if !result.Success {
		return decimal.Decimal{}, fmt.Errorf("result.Success - err")
	}

	return result.Result, nil
}
