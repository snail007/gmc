/*
 * @Description:Time tools
 * @Author: bearcat-panda
 * @Date: 2020-09-16 8:00:24
 * @LastEditTime: 2020-09-16 10:36:24
 * @LastEditors: bearcat-panda
 */
package timeutil

import (
	"fmt"
	"log"
	"strings"
	"time"
)

//Get the current time in seconds
func GetNowSecond() int64 {
	return time.Now().UTC().Unix()
}

//Get seconds based on time
func GetSecondByTime(t time.Time) int64 {
	return t.UTC().Unix()
}

//Get the current time in milliseconds
func GetNowMilliSecond() int64 {
	return time.Now().UTC().UnixNano() / 1e6
}

//Get the millisecond value based on time
func GetMilliSecondByTime(t time.Time) int64 {
	return t.UTC().UnixNano() / 1e6
}

//Change integer type to time type
func GetTimeByInt64(s int64) time.Time {
	length := len(fmt.Sprintf("%d", s)) //Get the length
	switch length {
	case len(fmt.Sprintf("%d", time.Second)): //second 10
		return time.Unix(s, 0)
	case len(fmt.Sprintf("%d", time.Second*1e3)): //millisecond 13
		return time.Unix(0, s*1e6)
	case len(fmt.Sprintf("%d", time.Second*1e9)): //nanosecond 19
		return time.Unix(0, s)
	default: //default
		panic("Length does not meet the requirements")
	}
	return time.Now() //currentTime
}

//Acquire the time according to the string
//isH true-24 hour clock false-12 hour clock
//No hours involved. The second field is ignored
func GetTimeByString(s string, isH bool) time.Time {
	value := strings.TrimSpace(s)
	if len(value) == 0 {
		panic("String format error")
	}

	var layout string
	switch isH {
	case true: //24h
		layout = "2006-01-02 15:04:05"
	case false: //12h
		layout = "2006-01-02 03:04:05"
	}

	if len(value) == 10 && strings.Contains(value, "-") {
		layout = "2006-01-02"
	}

	if len(value) == 10 && strings.Contains(value, "/") {
		layout = "2006/01/02"
	}

	if len(value) == 19 && strings.Contains(value, "/") && isH {
		layout = "2006/01/02 15:04:05"
	}

	if len(value) == 19 && strings.Contains(value, "/") && !isH {
		layout = "2006/01/02 03:04:05"
	}

	if len(value) == 14 && isH {
		layout = "20060102150405"
	}

	if len(value) == 14 && !isH {
		layout = "20060102030405"
	}

	if len(value) == 8 {
		layout = "20060102"
	}

	t, err := time.ParseInLocation(layout, value, time.Local)
	if err != nil {
		panic("String format error")
	}
	return t
}

//Get the zero point of incoming time
func GetZeroTime(t time.Time) time.Time {
	zeroTime, err := time.ParseInLocation("2006-01-02", t.Format("2006-01-02"), time.Local)
	if err != nil {
		panic("Time conversion error")
	}

	return zeroTime
}

//Get the zero time of the day
func GetNowZeroTime() time.Time {
	zeroTime, err := time.ParseInLocation(
		"2006-01-02",
		time.Now().Format("2006-01-02"),
		time.Local)
	if err != nil {
		panic("Time conversion error")
	}

	return zeroTime
}

//Get the first day of the month where the time is passed in.
// That is the zero point on the first day of a month
func GetFirstDateOfMonth(t time.Time) time.Time {
	t = t.AddDate(0, 0, -t.Day()+1)
	return GetZeroTime(t)
}

//Get the last day of the month where the time is passed in
//That is, 0 o'clock on the last day of a month
func GetLastDateOfMonth(t time.Time) time.Time {
	return GetFirstDateOfMonth(t).AddDate(0, 1, -1)
}

//Start of this year
func GetFirstDateOfYear(t time.Time) time.Time {
	t = t.AddDate(0, -int(t.Month())+1, -t.Day()+1)
	return GetZeroTime(t)
}

//End of this year
func GetLastDateOfYear(t time.Time) time.Time {
	return GetFirstDateOfYear(t).AddDate(1, 0, -1)
}

//Get the day of the week for the incoming time
func GetWeek(t time.Time) int {
	week := int(t.Weekday())
	if week == 0 {
		week = 7
	}
	return week
}

//Get the month of the year
func GetMonth(t time.Time) int {
	return int(t.Month())
}

//Is it a leap year
func IsLeapYear(t time.Time) bool {
	year := t.Year()
	if year%400 == 0 || (year%4 == 0 && year%100 != 0) {
		return true
	}
	return false
}

//Time zone conversion
func ParseWithLocation(location, timeStr string) (time.Time, error) {
	//Load time zone
	l, err := time.LoadLocation(location)
	if err != nil {
		log.Fatal(err)
		return time.Time{}, err
	}

	t, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, l)
	if err != nil {
		log.Fatal(err)
		return time.Time{}, err
	}

	return t, nil
}

//Convert to full time
func DateTimeToString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

//Convert to date
func DateToString(t time.Time) string {
	return t.Format("2006-01-02")
}
