package ticker

import "github.com/shopspring/decimal"

type Ticker interface {
	GetPrice(symbol string) (decimal.Decimal, error)
}
