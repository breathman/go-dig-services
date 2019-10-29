package timeseries

import (
	"github.com/breathman/go-dig-services/util"
)

type Service interface {
	PushTicker(market string, ticker Ticker) error
	GetTicker(market string) (Ticker, error)
	GetPriceCharts(market string, timeRange util.TimeRange) ([]Candle, error)
	Close()
}
