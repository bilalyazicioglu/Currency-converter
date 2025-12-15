package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

type ExchangeRateResponse struct {
	Result             string             `json:"result"`
	Documentation      string             `json:"documentation"`
	TermsOfUse         string             `json:"terms_of_use"`
	TimeLastUpdateUnix int64              `json:"time_last_update_unix"`
	TimeNextUpdateUnix int64              `json:"time_next_update_unix"`
	BaseCode           string             `json:"base_code"`
	ConversionRates    map[string]float64 `json:"conversion_rates"`
}

type ConversionResult struct {
	FromCurrency string
	ToCurrency   string
	Amount       float64
	Rate         float64
	Result       float64
	LastUpdated  time.Time
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://v6.exchangerate-api.com/v6",
	}
}

func (c *Client) GetExchangeRates(baseCurrency string) (*ExchangeRateResponse, error) {
	url := fmt.Sprintf("%s/%s/latest/%s", c.baseURL, c.apiKey, baseCurrency)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exchange rates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	var result ExchangeRateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Result != "success" {
		return nil, fmt.Errorf("API returned error: %s", result.Result)
	}

	return &result, nil
}

func (c *Client) Convert(fromCurrency, toCurrency string, amount float64) (*ConversionResult, error) {
	rates, err := c.GetExchangeRates(fromCurrency)
	if err != nil {
		return nil, err
	}

	rate, ok := rates.ConversionRates[toCurrency]
	if !ok {
		return nil, fmt.Errorf("currency %s not found", toCurrency)
	}

	return &ConversionResult{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Amount:       amount,
		Rate:         rate,
		Result:       amount * rate,
		LastUpdated:  time.Unix(rates.TimeLastUpdateUnix, 0),
	}, nil
}

func (c *Client) GetSupportedCurrencies() ([]string, error) {
	rates, err := c.GetExchangeRates("USD")
	if err != nil {
		return nil, err
	}

	currencies := make([]string, 0, len(rates.ConversionRates))
	for code := range rates.ConversionRates {
		currencies = append(currencies, code)
	}

	sort.Strings(currencies)
	return currencies, nil
}

func CommonCurrencies() []string {
	return []string{
		"USD", "EUR", "GBP", "JPY", "TRY",
		"CHF", "CAD", "AUD", "CNY", "INR",
		"KRW", "MXN", "BRL", "RUB", "ZAR",
	}
}

var CurrencyNames = map[string]string{
	"USD": "US Dollar",
	"EUR": "Euro",
	"GBP": "British Pound",
	"JPY": "Japanese Yen",
	"TRY": "Turkish Lira",
	"CHF": "Swiss Franc",
	"CAD": "Canadian Dollar",
	"AUD": "Australian Dollar",
	"CNY": "Chinese Yuan",
	"INR": "Indian Rupee",
	"KRW": "South Korean Won",
	"MXN": "Mexican Peso",
	"BRL": "Brazilian Real",
	"RUB": "Russian Ruble",
	"ZAR": "South African Rand",
	"AED": "UAE Dirham",
	"SAR": "Saudi Riyal",
	"SGD": "Singapore Dollar",
	"HKD": "Hong Kong Dollar",
	"NOK": "Norwegian Krone",
	"SEK": "Swedish Krona",
	"DKK": "Danish Krone",
	"PLN": "Polish Zloty",
	"THB": "Thai Baht",
	"IDR": "Indonesian Rupiah",
	"MYR": "Malaysian Ringgit",
	"PHP": "Philippine Peso",
	"CZK": "Czech Koruna",
	"ILS": "Israeli Shekel",
	"CLP": "Chilean Peso",
	"PKR": "Pakistani Rupee",
	"EGP": "Egyptian Pound",
	"TWD": "Taiwan Dollar",
	"VND": "Vietnamese Dong",
	"BDT": "Bangladeshi Taka",
	"ARS": "Argentine Peso",
	"COP": "Colombian Peso",
	"PEN": "Peruvian Sol",
	"UAH": "Ukrainian Hryvnia",
	"KZT": "Kazakhstani Tenge",
	"QAR": "Qatari Riyal",
	"KWD": "Kuwaiti Dinar",
	"BHD": "Bahraini Dinar",
	"OMR": "Omani Rial",
}
