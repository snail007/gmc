/*
 * @Description:Time tools
 * @Author: bearcat-panda
 * @Date: 2020-09-16 8:00:24
 * @LastEditTime: 2020-09-17 15:55:24
 * @LastEditors: bearcat-panda
 * godoc -http=:6060
 */
package timeutil

import (
	"fmt"
	"time"
)

//Get the current time in seconds
func ExampleUnix() {
	second := Unix()
	fmt.Println(second)
	//Output like : 1607702400
}

//Get seconds based on time
func ExampleTimeToUnix() {
	secondByTime := TimeToUnix(time.Now())
	fmt.Println(secondByTime)

	//Output like : 1607702400
}

//Get the current time in milliseconds
func ExampleMilliSeconds() {
	nowMilliSecond := MilliSeconds()
	fmt.Println(nowMilliSecond)

	//Output like : 1600320027565
}

//Get the millisecond value based on time
func ExampleTimeToMilliSeconds() {
	milliSecondByTime := TimeToMilliSeconds(time.Now())
	fmt.Println(milliSecondByTime)

	//Output like : 1600320027565
}

//Change integer type to time type
func ExampleInt64ToTime() {
	timeByInt := Int64ToTime(1600320027565) //milliSecond
	fmt.Println(timeByInt)

	timeByInt = Int64ToTime(1599546208) //second
	fmt.Println(timeByInt)

	//Output:
	//2020-09-17 13:20:27.565 +0800 CST
	//2020-09-08 14:23:28 +0800 CST
}

//Acquire the time according to the string
//@param s timString "2006-01-02","2006/01/02 15:04:05","2006/01/02 03:04:05"
//@param isH true-24 hour clock false-12 hour clock
//No hours involved. The second field is ignored
func ExampleStrToTime() {
	//No hours involved. The second field is ignored

	//true 2020-09-17 13:29:05
	timeByString := TimeFormat(MustStrToTime("2020-09-17"), "Y-m-d H:i:s")
	fmt.Println(timeByString)

	//true 2020-09-17 13:29:05
	timeByString = TimeFormat(MustStrToTime("2020-09-17 01:29:05pm"), "Y-m-d H:i:s")
	fmt.Println(timeByString)

	//false 2020-09-17 01:29:05
	timeByString = TimeFormat(MustStrToTime("2020-09-17 13:29:05"), "Y-m-d h:i:s")
	fmt.Println(timeByString)

	//Output:
	//2020-09-17 00:00:00
	//2020-09-17 13:29:05
	//2020-09-17 01:29:05
}

//Get the zero point of incoming time
func ExampleZeroTimeOf() {
	zeroTime := ZeroTimeOf(MustStrToTime("2020-09-17 00:00:00 +0800 CST"))
	fmt.Println(zeroTime)

	//Output:
	//2020-09-17 00:00:00 +0800 CST
}

//Get the zero time of the day
func ExampleZeroTime() {
	nowZeroTime := ZeroTime()
	fmt.Println(nowZeroTime)

	//Output like : 2020-09-17 00:00:00 +0800 CST
}

//Get the first day of the month where the time is passed in.
// That is the zero point on the first day of a month
func ExampleMonthZeroTimeOf() {
	firstDate := MonthZeroTimeOf(time.Now())
	fmt.Println(firstDate)

	//Output:
	//2020-09-01 00:00:00 +0800 CST
}

//Get the last day of the month where the time is passed in
//That is, 0 o'clock on the last day of a month
func ExampleMonthLastTimeOf() {
	lastDate := MonthLastTimeOf(time.Now())
	fmt.Println(lastDate)

	//Output:
	//2020-09-30 00:00:00 +0800 CST
}

//Start of this year
func ExampleYearZeroTimeOf() {
	firstData := YearZeroTimeOf(time.Now())
	fmt.Println(firstData)

	//Output:
	//2020-01-01 00:00:00 +0800 CST
}

//End of this year
func ExampleYearLastTimeOf() {
	lastDate := YearLastTimeOf(time.Now())
	fmt.Println(lastDate)

	//Output:
	//2020-12-31 00:00:00 +0800 CST
}

//Get the day of the week for the incoming time
func ExampleWeekdayOf() {
	week := WeekdayOf(MustStrToTime("2020-9-18 00:00:00"))
	fmt.Println(week)

	//Output:
	//5
}

//Get the month of the year
func ExampleMonthOf() {
	month := MonthOf(MustStrToTime("2020-9-18 00:00:00"))
	fmt.Println(month)

	//Output:
	//9
}

//Is it a leap year
func ExampleIsLeapYear() {
	flag := IsLeapYear(MustStrToTime("2020/09/16")) //ture
	fmt.Println(flag)

	flag = IsLeapYear(MustStrToTime("2007/09/16")) //false
	fmt.Println(flag)

	//Output:
	//true
	//false
}

//Time zone conversion
//@param location "","UTC","America/New_York"
//@param timeStrEx	"2006-01-02 15:04:05"
func ExampleStrToTimeOfLocation() {
	t, err := StrToTimeOfLocation("America/New_York", "2020-09-17 13:44:05")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(t)

	//Output:
	//2020-09-17 13:44:05 -0400 EDT
}

//Convert to full time
func ExampleTimeToStr() {
	s := TimeToStr(MustStrToTime("2007/09/16"))
	fmt.Println(s)

	//Output:
	//2007-09-16 00:00:00
}

//Convert to date
func ExampleTimeToDateStr() {
	s := TimeToDateStr(MustStrToTime("2007/09/16"))
	fmt.Println(s)

	//Output:
	//2007-09-16
}

func ExampleTimeFormat() {
	now := MustStrToTime("2020/09/17 01:50:04")
	s := TimeFormat(now, "h:i:s")
	fmt.Println(s)

	s = TimeFormat(now, "Y/m/d")
	fmt.Println(s)

	s = TimeFormat(now, "Y-m-d h:i:s")
	fmt.Println(s)

	s = TimeFormat(now, "yyyy-mm-dd hh:ii:ss")
	fmt.Println(s)

	/*
		Y-2006	y-06
		m-01 	n-1
		d-02 	j-2
		H-15	h-03	g-3
		i-04
		s-05
	*/

	//Output:
	//01:50:04
	//2020/09/17
	//2020-09-17 01:50:04
	//2020-09-17 01:50:04

}
