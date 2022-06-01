package rest_api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

const (
	timeout = time.Second * 10
)

type Client struct {
	client   http.Client
	apiKey   string
	basePath string
}

type (
	Response struct {
		Date       string          `json:"date"`
		Historical string          `json:"historical"`
		Info       Info            `json:"info"`
		Query      Query           `json:"query"`
		Result     decimal.Decimal `json:"result"`
		Success    bool            `json:"success"`
	}
	Query struct {
		Amount decimal.Decimal `json:"amount"`
		From   string          `json:"from"`
		To     string          `json:"to"`
	}
	Info struct {
		Rate      decimal.Decimal `json:"rate"`
		Timestamp int             `json:"timestamp"`
	}
)

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

func (c *Client) ExchangeCurrency(ctx context.Context, amount decimal.Decimal, currency string) (decimal.Decimal, error) {
	json.Unmarshal()
}

func (c *Client) request(ctx context.Context, amount decimal.Decimal, currency string) (decimal.Decimal, error) {

}
