package rest_api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/shopspring/decimal"
)

const (
	timeout         = time.Second * 10
	convertURL      = "/convert"
	defaultCurrency = "RUB"
)

type Response struct {
	Result  decimal.Decimal `json:"result"`
	Success bool            `json:"success"`
}

var (
	errUnknown = errors.New("unknown")
)

func (c *Client) ExchangeCurrency(ctx context.Context, amount decimal.Decimal, toCurrency string) (decimal.Decimal, error) {
	url, err := c.ConvertToUrl(amount, toCurrency)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("c.ConvertToUrl")
	}

	result, err := c.request(ctx, url, http.MethodGet, nil)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("c.request: %w", err)
	}

	return result, nil
}

func (c *Client) ConvertToUrl(money decimal.Decimal, currency string) (string, error) {
	endPoint := fmt.Sprintf("%s%s", c.basePath, convertURL)

	urlAddres, err := url.Parse(endPoint)
	if err != nil {
		return "", fmt.Errorf("url.Parse: %w", err)
	}

	param := url.Values{}

	param.Add("apikey", c.apiKey)
	param.Add("from", defaultCurrency)
	param.Add("to", currency)
	param.Add("amount", money.String())

	urlAddres.RawQuery = param.Encode()

	return urlAddres.String(), nil

}
