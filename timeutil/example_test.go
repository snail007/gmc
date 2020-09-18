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
func ExampleGetNowSecond() {
	second := GetNowSecond()
	fmt.Println(second)

	//Output:
	//1599546208
}

//Get seconds based on time
func ExampleGetSecondByTime() {
	secondByTime := GetSecondByTime(time.Now())
	fmt.Println(secondByTime)

	//Output:
	//1599546208
}

//Get the current time in milliseconds
func ExampleGetNowMilliSecond() {
	nowMilliSecond := GetNowMilliSecond()
	fmt.Println(nowMilliSecond)

	//Output:
	//1600320027565
}

//Get the millisecond value based on time
func ExampleGetMilliSecondByTime() {
	milliSecondByTime := GetMilliSecondByTime(time.Now())
	fmt.Println(milliSecondByTime)

	//Output:
	//1600320027565
}

//Change integer type to time type
func ExampleGetTimeByInt64() {
	timeByInt := GetTimeByInt64(1600320027565) //milliSecond
	fmt.Println(timeByInt)

	timeByInt = GetTimeByInt64(1599546208) //second
	fmt.Println(timeByInt)

	//Output:
	//2020-09-17 13:20:27.565 +0800 CST
	//2020-09-08 14:23:28 +0800 CST
}

//Acquire the time according to the string
//@param s timString "2006-01-02","2006/01/02 15:04:05","2006/01/02 03:04:05"
//@param isH true-24 hour clock false-12 hour clock
//No hours involved. The second field is ignored
func ExampleGetTimeByString() {
	//No hours involved. The second field is ignored
	timeByString := GetTimeByString("2020-09-17", false)
	fmt.Println(timeByString)

	//true 2020-09-17 13:29:05
	//false 2020-09-17 01:29:05
	timeByString = GetTimeByString("2020-09-17 13:29:05", true)
	fmt.Println(timeByString)

	//Output:
	//true 2020-09-17 13:29:05
	//false 2020-09-17 01:29:05
}

//Get the zero point of incoming time
func ExampleGetZeroTime() {
	zeroTime := GetZeroTime(time.Now())
	fmt.Println(zeroTime)

	//Output:
	//2020-09-17 00:00:00 +0800 CST
}

//Get the zero time of the day
func ExampleGetNowZeroTime() {
	nowZeroTime := GetNowZeroTime()
	fmt.Println(nowZeroTime)

	//Output:
	//2020-09-17 00:00:00 +0800 CST
}

//Get the first day of the month where the time is passed in.
// That is the zero point on the first day of a month
func ExampleGetFirstDateOfMonth() {
	firstDate := GetFirstDateOfMonth(time.Now())
	fmt.Println(firstDate)

	//Output:
	//2020-09-01 00:00:00 +0800 CST
}

//Get the last day of the month where the time is passed in
//That is, 0 o'clock on the last day of a month
func ExampleGetLastDateOfMonth() {
	lastDate := GetLastDateOfMonth(time.Now())
	fmt.Println(lastDate)

	//Output:
	//2020-09-30 00:00:00 +0800 CST
}

//Start of this year
func ExampleGetFirstDateOfYear() {
	firstData := GetFirstDateOfYear(time.Now())
	fmt.Println(firstData)

	//Output:
	//2020-01-01 00:00:00 +0800 CST
}

//End of this year
func ExampleGetLastDateOfYear() {
	lastDate := GetLastDateOfYear(time.Now())
	fmt.Println(lastDate)

	//Output:
	//2020-12-31 00:00:00 +0800 CST
}

//Get the day of the week for the incoming time
func ExampleGetWeek() {
	week := GetWeek(time.Now())
	fmt.Println(week)

	//Output:
	//2020-9-17 | Thursday | 4
}

//Get the month of the year
func ExampleGetMonth() {
	month := GetMonth(time.Now())
	fmt.Println(month)

	//Output:
	//2020-9-17 | 9
}

//Is it a leap year
func ExampleIsLeapYear() {
	flag := IsLeapYear(GetTimeByString("2020/09/16", false)) //ture
	fmt.Println(flag)

	flag = IsLeapYear(GetTimeByString("2007/09/16", false)) //false
	fmt.Println(flag)

	//Output:
	//2020/9/16 true | 2007/9/16 false
}

//Time zone conversion
//@param location "","UTC","America/New_York"
//@param timeStrEx	"2006-01-02 15:04:05"
func ExampleParseWithLocation() {
	t, err := ParseWithLocation("America/New_York", "2020-09-17 13:44:05")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(t)

	//Output:
	//2020-09-17 13:44:05 -0400 EDT
}

//Convert to full time
func ExampleDateTimeToString() {
	s := DateTimeToString(time.Now())
	fmt.Println(s)

	//Output:
	//2020-09-17 13:47:07
}

//Convert to date
func ExampleDateToString() {
	s := DateToString(time.Now())
	fmt.Println(s)

	//Output:
	//2020-09-17
}

func ExampleTimeFormat() {
	s := TimeFormat(time.Now(), "h:i:s")
	fmt.Println(s)

	s = TimeFormat(time.Now(), "Y/m/d")
	fmt.Println(s)

	s = TimeFormat(time.Now(), "Y-m-d h:i:s")
	fmt.Println(s)

	s = TimeFormat(time.Now(), "yyyy-mm-dd hh:ii:ss")
	fmt.Println(s)

	//Output:
	//"h:i:s" 01:50:04
	//"Y/m/d" 2020/09/17
	//"Y-m-d h:i:s" 2020-09-17 01:50:04
	//"yyyy-mm-dd hh:ii:ss"  20-09-17 01:50:04
	/** 2006-01-02 03:04:05 |  2006-01-02 15:04:05
	Y-2006	y-06
	m-01 	n-1
	d-02 	j-2
	H-15	h-03	g-3
	i-04
	s-05
	*/
}
