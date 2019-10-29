package timeseries

import (
	"time"

	"github.com/shopspring/decimal"
)

type Candle struct {
	Open  decimal.Decimal `json:"open"`
	Close decimal.Decimal `json:"close"`
	Low   decimal.Decimal `json:"low"`
	High  decimal.Decimal `json:"high"`
	DT    time.Time       `json:"dt"`
}

type Ticker struct {
	Price decimal.Decimal `json:"price"`
	DT    time.Time       `json:"datetime"`
}
