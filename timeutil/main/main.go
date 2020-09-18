package main

import (
	"github.com/snail007/gmc/timeutil"
	"log"
	"time"
)

func main() {
	//Get the current time in seconds
	second := timeutil.Unix()
	log.Println("second", second)

	//Get seconds based on time
	secondByTime := timeutil.TimeToUnix(time.Now())
	log.Println("secondByTime", secondByTime)

	//Get the current time in milliseconds
	milliSecond := timeutil.MilliSeconds()
	log.Println("milliSecond", milliSecond)

	//Get the millisecond value based on time
	milliSecondByTime := timeutil.TimeToMilliSeconds(time.Now())
	log.Println("milliSecondByTime", milliSecondByTime)

	//Change integer type to time type
	timeByInt := timeutil.Int64ToTime(time.Now().Unix())
	log.Println("timeByInt", timeByInt)

	//Acquire the time according to the string
	//isH true-24 hour clock false-12 hour clock
	//No hours involved. The second field is ignored
	timeByString1 := timeutil.StrToTime("2020-09-16", false)
	timeByString2 := timeutil.StrToTime("2020-09-16 10:27:00", true)
	log.Println("timeByString1", timeByString1)
	log.Println("timeByString2", timeByString2)

	//Get the zero point of incoming time
	zeroTime := timeutil.ZeroTimeOf(time.Now())
	log.Println("zeroTime", zeroTime)

	//Get the zero time of the day
	nowZeroTime := timeutil.ZeroTime()
	log.Println("nowZeroTime", nowZeroTime)

	//Get the first day of the month where the time is passed in.
	// That is the zero point on the first day of a month
	firstDateOfMonth := timeutil.MonthZeroTimeOf(time.Now())
	log.Println("firstDateOfMonth", firstDateOfMonth)

	//Get the last day of the month where the time is passed in
	//That is, 0 o'clock on the last day of a month
	lastDateOfMonth := timeutil.MonthLastTimeOf(time.Now())
	log.Println("lastDateOfMonth", lastDateOfMonth)

	//Start of this year
	firstDateOfYear := timeutil.YearZeroTimeOf(time.Now())
	log.Println("firstDateOfYear", firstDateOfYear)

	//End of this year
	lastDateOfYear := timeutil.YearLastTimeOf(time.Now())
	log.Println("lastDateOfYear", lastDateOfYear)

	//Get the day of the week for the incoming time
	week := timeutil.WeekdayOf(time.Now())
	log.Println("week", week)

	//Get the month of the year
	month := timeutil.MonthOf(time.Now())
	log.Println("month", month)

	//Is it a leap year
	isLeapYear := timeutil.IsLeapYear(time.Now())
	log.Println("isLeapYear", isLeapYear)

	//Time zone conversion
	parseLocate, err := timeutil.StrToTimeOfLocation("UTC", "2020-09-17 10:23:01")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("parseLocate", parseLocate)

	//Convert to full time
	timeToString := timeutil.TimeToStr(time.Now())
	log.Println("timeToString", timeToString)

	//Convert to date
	dateToString := timeutil.TimeToDateStr(time.Now())
	log.Println("dateToString", dateToString)

	//Time format
	timeFormat1 := timeutil.TimeFormat(time.Now(), "yyyy-dd-mm")
	timeFormat2 := timeutil.TimeFormat(time.Now(), "yyyy-dd-mm hh:ii:ss")
	timeFormat3 := timeutil.TimeFormat(time.Now(), "Y-d-m h:i:s")
	log.Println("timeFormat1", timeFormat1)
	log.Println("timeFormat2", timeFormat2)
	log.Println("timeFormat3", timeFormat3)
}
