package timeseries

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/shopspring/decimal"

	"github.com/breathman/go-dig-services/config"
	"github.com/breathman/go-dig-services/log"
	"github.com/breathman/go-dig-services/util"
)

type InfluxTS struct {
	client influx.Client
	dbName string
	log    *log.CtxLogger
}

func NewInfluxTS(conf *config.APPConfig, log *log.Service) (Service, error) {
	ctxLog := log.NewPrefix("influx")

	client, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr: "http://" + conf.Influx.Host + ":" + conf.Influx.Port,
	})
	return &InfluxTS{
		client: client,
		dbName: conf.Influx.DBName,
		log:    ctxLog,
	}, err
}

func (its *InfluxTS) PushTicker(market string, ticker Ticker) error {
	tags := map[string]string{"market": market}
	price, _ := ticker.Price.Float64()
	fields := map[string]interface{}{
		"price": price,
	}
	return its.addMeasurement("ticker", tags, fields)
}

func (its *InfluxTS) GetTicker(market string) (Ticker, error) {
	query := fmt.Sprintf("SELECT LAST(price) FROM ticker WHERE market='%s'", market) //nolint
	values, err := its.query(query)
	if err != nil {
		return Ticker{}, err
	}
	if len(values) == 0 {
		return Ticker{}, fmt.Errorf("ticker not found")
	}
	row := values[0]
	if len(row) != 2 {
		return Ticker{}, fmt.Errorf("ticker extracting error with values: %v", row)
	}
	ticker := Ticker{}
	ticker.DT, err = time.Parse(time.RFC3339, row[0].(string))
	if err != nil {
		return Ticker{}, err
	}
	ticker.Price, err = getDecimal(row[1])
	if err != nil {
		return Ticker{}, err
	}
	return ticker, nil
}

func (its *InfluxTS) GetPriceCharts(market string, timeRange util.TimeRange) ([]Candle, error) {
	query := fmt.Sprintf("SELECT FIRST(price), LAST(price), MIN(price), MAX(price) FROM ticker WHERE market='%s' AND time >= '%s' AND time <= '%s' GROUP BY time(%s) fill(none) ORDER BY DESC tz('Europe/Moscow')", market, timeRange.Start.Format(time.RFC3339), timeRange.End.Format(time.RFC3339), getSpan(timeRange)) // nolint
	values, err := its.query(query)
	if err != nil {
		return []Candle{}, err
	}
	candles := make([]Candle, len(values))
	for i, row := range values {
		candle := Candle{}
		candle.DT, err = time.Parse(time.RFC3339, row[0].(string))
		if err != nil {
			return []Candle{}, err
		}
		if candle.Open, err = getDecimal(row[1]); err != nil {
			return []Candle{}, err
		}
		if candle.Close, err = getDecimal(row[2]); err != nil {
			return []Candle{}, err
		}
		if candle.Low, err = getDecimal(row[3]); err != nil {
			return []Candle{}, err
		}
		if candle.High, err = getDecimal(row[4]); err != nil {
			return []Candle{}, err
		}
		candles[i] = candle
	}
	return candles, nil
}

func getDecimal(value interface{}) (decimal.Decimal, error) {
	res, err := decimal.NewFromString(value.(json.Number).String())
	if err != nil {
		return res, err
	}
	return res, nil
}

func getSpan(timeRange util.TimeRange) string {
	dur := timeRange.End.Sub(timeRange.Start).Hours()

	// group by time period
	switch {
	case dur < 4:
		return "15m" // 3h to 15m
	case dur < 25 && dur > 4:
		return "1h" // 1d to 1h
	case dur < 169 && dur > 25:
		return "6h" // 1w to 6h
	case dur < 755 && dur > 169:
		return "1d" // 1m to 1d
	case dur < 2263 && dur > 755:
		return "3d" // 3m to 3d
	case dur < 4525 && dur > 2263:
		return "6d" // 6m to 6d
	case dur < 9048 && dur > 4525:
		return "12d" // 1y to 1w
	case dur > 9048:
		return "4w" // all time to 1m
	default:
		return "1d"
	}
}

func (its *InfluxTS) addMeasurement(measurementName string, tags map[string]string, fields map[string]interface{}) error {
	its.log.Debugf("add new %q measurement", measurementName)
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database: its.dbName,
	})
	if err != nil {
		return err
	}

	point, err := influx.NewPoint(measurementName, tags, fields, time.Now())
	if err != nil {
		return err
	}

	bp.AddPoint(point)

	// Write the batch
	if err := its.client.Write(bp); err != nil {
		return err
	}
	its.log.Debugf("%q measurement added", measurementName)

	return nil
}

func (its *InfluxTS) query(query string) ([][]interface{}, error) {
	defer its.recovery()

	its.log.Debugf("query to Influx: %q", query)
	q := influx.Query{
		Command:  query,
		Database: its.dbName,
	}
	response, err := its.client.Query(q)
	if err != nil {
		return nil, err
	}
	if response.Error() != nil {
		return nil, response.Error()
	}

	result := response.Results
	if len(result) < 1 || len(result[0].Series) < 1 {
		err = errors.New("time series its not detected")
		return nil, err
	}
	return result[0].Series[0].Values, nil
}

func (its *InfluxTS) Close() {
	if err := its.client.Close(); err != nil {
		its.log.Errorf("Close connection error %q", err.Error())
	}
}

func (its *InfluxTS) recovery() {
	if p := recover(); p != nil {
		its.log.Errorf("Recovered in %q", p)
	}
}
