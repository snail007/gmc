package gutil

import (
	"strings"
	"time"
)

func location(loc ...*time.Location) *time.Location {
	var loc0 *time.Location
	if len(loc) == 1 {
		loc0 = loc[0]
	} else {
		loc0 = time.UTC
	}
	return loc0
}

func getTime(t ...time.Time) time.Time {
	var t0 time.Time
	if len(t) == 1 {
		t0 = t[0]
	} else {
		t0 = time.Now()
	}
	return t0
}

//Get the millisecond value based on time
func MilliSecondsOf(t ...time.Time) int64 {
	return getTime(t...).UnixNano() / 1e6
}

//Get the zero point of incoming time
func ZeroMinuteOf(t_ ...time.Time) time.Time {
	t := getTime(t_...)
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
}

//Get the zero point of incoming time
func ZeroHourOf(t_ ...time.Time) time.Time {
	t := getTime(t_...)
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
}

//Get the zero point of incoming time
func ZeroDayOf(t_ ...time.Time) time.Time {
	t := getTime(t_...)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

//Get the zero point of incoming time
func ZeroMonthOf(t_ ...time.Time) time.Time {
	t := getTime(t_...)
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

//Get the zero point of incoming time
func ZeroYearOf(t_ ...time.Time) time.Time {
	t := getTime(t_...)
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

//Get the day of the week for the incoming time
func WeekdayOf(t_ ...time.Time) int {
	t := getTime(t_...)
	week := int(t.Weekday())
	if week == 0 {
		week = 7
	}
	return week
}

//Get the month of the year
func MonthOf(t_ ...time.Time) int {
	t := getTime(t_...)
	return int(t.Month())
}

func parseAny(datestr string, loc_ ...interface{}) (t time.Time, err error) {
	var loc0 interface{}
	var loc *time.Location
	if len(loc_) == 1 {
		loc0 = loc_[0]
	}
	switch v := loc0.(type) {
	case *time.Location:
		loc = v
	case string:
		loc, err = time.LoadLocation(v)
	default:
		loc = time.UTC
	}
	if err != nil {
		return time.Time{}, err
	}
	p, err := parseTime(datestr, loc)
	if err != nil {
		return time.Time{}, err
	}
	return p.parse()
}

//Time format text
//@param format  "%h:%i:%s","%Y-%m-%d %h:%i:%s","%y-%m-%d"
func TimeFormatText(t time.Time, text string) string {
	d := map[string]string{
		"Y": t.Format("2006"),
		"y": t.Format("06"),
		"m": t.Format("01"),
		"n": t.Format("1"),
		"d": t.Format("02"),
		"j": t.Format("2"),
		"H": t.Format("15"),
		"h": t.Format("03"),
		"g": t.Format("3"),
		"i": t.Format("04"),
		"s": t.Format("05"),
	}
	for k, v := range d {
		k = "%" + k
		text = strings.Replace(text, k, v, 1)
	}
	return text
}
