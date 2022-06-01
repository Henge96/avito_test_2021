package rest_api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

func Convert(money decimal.Decimal, currency string) (decimal.Decimal, error) {

	var (
		response = Response{}
		url      = "https://api.apilayer.com/exchangerates_data/convert?to=RUB"
		total    = "&amount=" + money.StringFixedBank(4)
		from     = "&from=" + currency
	)

	client := &http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest("GET", url+from+total, nil)
	if err != nil {
		return decimal.Decimal{}, err
	}
	req.Header.Set("apikey", "Aa5vJTXRCS5gXJinXKFoJUlFKCqtaSTp")
	result, err := client.Do(req)
	if err != nil {
		return decimal.Decimal{}, err
	}
	defer result.Body.Close()

	err = json.NewDecoder(result.Body).Decode(&response)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return response.Result, nil

}
