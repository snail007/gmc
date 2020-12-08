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
	"regexp"
	"strings"
	"time"
)

//Get the current time in seconds
func Unix() int64 {
	return time.Now().UTC().Unix()
}

//Get seconds based on time
func TimeToUnix(t time.Time) int64 {
	return t.UTC().Unix()
}

//Get the current time in milliseconds
func MilliSeconds() int64 {
	return time.Now().UTC().UnixNano() / 1e6
}

//Get the millisecond value based on time
func TimeToMilliSeconds(t time.Time) int64 {
	return t.UTC().UnixNano() / 1e6
}

// Int64ToTime converts integer type to time type
func Int64ToTime(s int64) (r time.Time) {
	defer func() {
		r.In(time.Local)
	}()
	length := len(fmt.Sprintf("%d", s)) //Get the length
	switch length {
	case len(fmt.Sprintf("%d", time.Second)): //second 10
		return time.Unix(s, 0)
	case len(fmt.Sprintf("%d", time.Second*1e3)): //millisecond 13
		return time.Unix(0, s*1e6)
	case len(fmt.Sprintf("%d", time.Second*1e9)): //nanosecond 19
		return time.Unix(0, s)
	default: //default
		return time.Time{}
	}
	return time.Now() //currentTime
}

// StrToTime with Location, equivalent to time.ParseInLocation() timezone/offset
// rules.  Using location arg, if timezone/offset info exists in the
// datestring, it uses the given location rules for any zone interpretation.
// That is, MST means one thing when using America/Denver and something else
// in other locations.
func StrToTime(s string) (time.Time, error) {
	return ParseAny(s)
}

// MustStrToTime  parse a date, and zero time.Time{} returned if it can't be parsed.  Used for testing.
// Not recommended for most use-cases.
func MustStrToTime(s string) (r time.Time) {
	r, _ = ParseAny(s)
	r.In(time.Local)
	return
}

// StrToTimeIn with Location, equivalent to time.ParseInLocation() timezone/offset
// rules.  Using location arg, if timezone/offset info exists in the
// datestring, it uses the given location rules for any zone interpretation.
// That is, MST means one thing when using America/Denver and something else
// in other locations.
func StrToTimeIn(s string, loc *time.Location) (time.Time, error) {
	return parseIn(s, loc)
}

// MustStrToTime  parse a date, and zero time.Time{} returned if it can't be parsed.  Used for testing.
// Not recommended for most use-cases.
func MustStrToTimeIn(s string, loc *time.Location) (t time.Time) {
	var err error
	t, err = parseIn(s, loc)
	if err != nil {
		t = time.Time{}
		return
	}
	return
}

//Get the zero point of incoming time
func ZeroTimeOf(t time.Time) time.Time {
	return time.Date(t.Year(),t.Month(),t.Day(),0,0,0,0,t.Location())
}

//Get the zero time of the day
func ZeroTime() time.Time {
	zeroTime, err := time.ParseInLocation(
		"2006-01-02",
		time.Now().Format("2006-01-02"),
		time.Local)
	if err != nil {
		return time.Time{}.In(time.Local)
	}

	return zeroTime
}

//Get the first day of the month where the time is passed in.
// That is the zero point on the first day of a month
func MonthZeroTimeOf(t time.Time) time.Time {
	t = t.AddDate(0, 0, -t.Day()+1)
	return ZeroTimeOf(t)
}

//Get the last day of the month where the time is passed in
//That is, 0 o'clock on the last day of a month
func MonthLastTimeOf(t time.Time) time.Time {
	return MonthZeroTimeOf(t).AddDate(0, 1, -1)
}

//Start of this year
func YearZeroTimeOf(t time.Time) time.Time {
	t = t.AddDate(0, -int(t.Month())+1, -t.Day()+1)
	return ZeroTimeOf(t)
}

//End of this year
func YearLastTimeOf(t time.Time) time.Time {
	return YearZeroTimeOf(t).AddDate(1, 0, -1)
}

//Get the day of the week for the incoming time
func WeekdayOf(t time.Time) int {
	week := int(t.Weekday())
	if week == 0 {
		week = 7
	}
	return week
}

//Get the month of the year
func MonthOf(t time.Time) int {
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
//@param location "","UTC","America/New_York"
//@param timeStrEx	"2006-01-02 15:04:05"
func StrToTimeOfLocation(location, timeStr string) (time.Time, error) {
	//Load time zone
	l, err := time.LoadLocation(location)
	if err != nil {
		// log.Println(err)
		return time.Time{}.In(time.Local), err
	}

	t, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, l)
	if err != nil {
		// log.Println(err)
		return time.Time{}.In(time.Local), err
	}

	return t, nil
}

//Convert to full time
func TimeToStr(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

//Convert to date
func TimeToDateStr(t time.Time) string {
	return t.Format("2006-01-02")
}

//Time format
//@param format  "h:i:s","Y-m-d h:i:s","y-m-d"
func TimeFormat(t time.Time, format string) string {
	format = strings.TrimSpace(format)
	format = strings.Replace(format, "yyyy", "Y", 1)
	format = stringDedup(format) //String deduplication
	d := map[string]string{
		"Y": "2006",
		"y": "06",
		"m": "01",
		"n": "1",
		"d": "02",
		"j": "2",
		"H": "15",
		"h": "03",
		"g": "3",
		"i": "04",
		"s": "05",
	}
	for k, v := range d {
		format = strings.Replace(format, k, v, 1)
	}
	return t.Format(format)
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

//String deduplication
func stringDedup(format string) string {
	var maps = make(map[string]int)
	var temp []int32
	var index int
	for _, v := range format {
		if _, ok := maps[string(v)]; !ok &&
			regexp.MustCompile("[a-zA-Z]").MatchString(string(v)) { //Add letters
			maps[string(v)] = index
			temp = append(temp, v)
			index++
		}
		if regexp.MustCompile("[^a-zA-Z]").MatchString(string(v)) { //Add non-letter
			temp = append(temp, v)
			index++
		}
	}
	return string(temp)
}

//-----------------------------------------
