package util

import (
	"time"

	"github.com/pkg/errors"
)

type TimeRange struct {
	Start time.Time
	End   time.Time
	Loc   time.Location
}

func SetTimeShort(t time.Time, hour int, min int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), hour, min, 0, 0, t.Location())
}

func (r *TimeRange) ContainsInclusive(t time.Time) bool {
	return t.Equal(r.Start) || (t.After(r.Start) && t.Before(r.End)) || t.Equal(r.End)
}

func NewTimeRangeFromStrings(start, end string, loc *time.Location) (TimeRange, error) {
	var (
		timeRange TimeRange
		err       error
	)

	timeRange.Start, err = BuildTimeFromString(start, loc)
	if err != nil {
		return TimeRange{}, errors.Errorf("cannot parse start date: %q", err)
	}
	timeRange.End, err = BuildTimeFromString(end, loc)
	if err != nil {
		return TimeRange{}, errors.Errorf("cannot parse end date: %q", err)
	}
	timeRange.Loc = *loc

	return timeRange, nil
}

func BuildTimeFromString(strTime string, loc *time.Location) (time.Time, error) {
	var (
		res time.Time
		err error
	)

	if loc != nil {
		res, err = time.ParseInLocation("2006-01-02T15:04:05", strTime, loc)
		if err != nil {
			return res, err
		}
	} else {
		res, _ = time.ParseInLocation("2006-01-02T15:04:05", strTime, time.UTC)
	}

	return res, nil
}

func TimeInRange(t time.Time, start time.Time, end time.Time) bool {
	return (t.After(start) || t.Equal(start)) && (t.Before(end) || t.Equal(end))
}
