package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/handlename/go-mysql-slowlog-parser"
)

func initRange(begin string, end string) (string, string) {
	b := begin
	if begin == "" {
		b = "0"
	}

	e := end
	if end == "" {
		e = "0"
	}

	return b, e
}

func initRangeTime(begin string, end string) (string, string) {
	b := begin
	if begin == "" {
		b = "0000-01-01T00:00:00"
	}

	e := end
	if end == "" {
		e = "0000-01-01T00:00:00"
	}

	return b, e
}

func float32Range(val float32, begin, end string) bool {
	b, e := initRange(begin, end)
	beginf, err := strconv.ParseFloat(b, 32)
	if err != nil {
		return false
	}

	endf, err := strconv.ParseFloat(e, 32)
	if err != nil {
		return false
	}

	if beginf == 0 && endf != 0 {
		return val <= float32(endf)
	}

	if beginf != 0 && endf == 0 {
		return val > float32(beginf)
	}

	return val >= float32(beginf) && val < float32(endf)
}

func int32Range(val int32, begin, end string) bool {
	b, e := initRange(begin, end)
	begini, err := strconv.ParseInt(b, 10, 32)
	if err != nil {
		return false
	}

	endi, err := strconv.ParseInt(e, 10, 32)
	if err != nil {
		return false
	}

	if begini == 0 && endi != 0 {
		return val <= int32(endi)
	}

	if begini != 0 && endi == 0 {
		return val > int32(begini)
	}

	return val >= int32(begini) && val < int32(endi)
}

func offsetToTime(offset int) string {
	if offset == 0 {
		return fmt.Sprintf("+00:00")
	}

	var format string
	hour := offset / 3600
	min := (offset / 60) % 60

	if hour < 0 {
		hour *= -1
	}

	if min < 0 {
		min *= -1
	}

	if offset < 0 {
		format = "-%02d:%02d"
	} else {
		format = "+%02d:%02d"
	}

	return fmt.Sprintf(format, hour, min)
}

func parseTime(ts string, loc *time.Location) (time.Time, error) {
	if ts == "0000-01-01T00:00:00" {
		var t time.Time
		return t, nil
	}

	t, err := time.Parse("15:04:05", ts)
	_, offset := time.Now().In(time.Local).Zone()
	offsetTime := offsetToTime(offset)

	if err == nil {
		now := time.Now()

		t, err = time.Parse("2006-01-02T15:04:05-07:00", fmt.Sprintf("%d-%02d-%02dT%s%s", now.Year(), now.Month(), now.Day(), ts, offsetTime))

		return t.In(loc), err
	}

	t, err = time.Parse("2006-01-02T15:04:05-07:00", fmt.Sprintf("%s%s", ts, offsetTime))

	return t, err
}

func timeRange(val int64, begin, end string, loc *time.Location) bool {
	b, e := initRangeTime(begin, end)

	beginT, err := parseTime(b, loc)
	if err != nil {
		return false
	}

	endT, err := parseTime(e, loc)
	if err != nil {
		return false
	}

	if !beginT.IsZero() && endT.IsZero() {
		return val >= beginT.Unix()
	}

	if beginT.IsZero() && !endT.IsZero() {
		return val < endT.In(loc).Unix()
	}

	return val >= beginT.Unix() && val < endT.Unix()
}

func VerifyRange(s slowlog.Parsed, c Config) bool {
	bools := make([]bool, 0)

	if c.QueryTimeBegin != "" || c.QueryTimeEnd != "" {
		bools = append(bools, float32Range(s.QueryTime, c.QueryTimeBegin, c.QueryTimeEnd))
	}

	if c.LockTimeBegin != "" || c.LockTimeEnd != "" {
		bools = append(bools, float32Range(s.LockTime, c.LockTimeBegin, c.LockTimeEnd))
	}

	if c.RowsSentBegin != "" || c.RowsSentEnd != "" {
		bools = append(bools, int32Range(s.RowsSent, c.RowsSentBegin, c.RowsSentEnd))
	}

	if c.RowsExaminedBegin != "" || c.RowsExaminedEnd != "" {
		bools = append(bools, int32Range(s.RowsExamined, c.RowsExaminedBegin, c.RowsExaminedEnd))
	}

	if c.TimeBegin != "" || c.TimeEnd != "" {
		bools = append(bools, timeRange(s.Datetime, c.TimeBegin, c.TimeEnd, c.Location))
	}

	// or
	if c.RangeOr {
		for _, b := range bools {
			if b {
				return true
			}
		}

		return false
	}

	// and
	for _, b := range bools {
		if !b {
			return false
		}
	}

	return true
}
