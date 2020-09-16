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
	"testing"
	"time"
)

//Get the current time in seconds
func Test_GetNowSecond(t *testing.T) {
	nowSecond := GetNowSecond()
	t.Log(nowSecond)
}

//Get seconds based on time
func Test_GetSecondByTime(t *testing.T) {
	second := GetSecondByTime(time.Now())
	t.Log(second)
}

//Get the current time in milliseconds
func Test_GetNowMilliSecond(t *testing.T) {
	nowMilliSecond := GetNowMilliSecond()
	t.Log(nowMilliSecond)
}

//Get the millisecond value based on time
func Test_GetMilliSecondByTime(t *testing.T) {
	milliSecond := GetMilliSecondByTime(time.Now())
	t.Log(milliSecond)
}

//Change integer type to time type
func Test_GetTimeByInt64(t *testing.T) {
	second := GetNowSecond()            //Get the second value of the current time
	milliSecond := GetNowMilliSecond()  //Get the current time in milliseconds
	nanoSecond := time.Now().UnixNano() //Get the nanosecond value of the current time

	t.Log(len(fmt.Sprintf("%d", second)))      //The length of the second value
	t.Log(len(fmt.Sprintf("%d", milliSecond))) //The length of the millisecond value
	t.Log(len(fmt.Sprintf("%d", nanoSecond)))  //The length of the nanosecond value

	//Convert to time type
	secondTime := GetTimeByInt64(second)
	milliTime := GetTimeByInt64(milliSecond)
	nanoTime := GetTimeByInt64(int64(nanoSecond))

	//print
	t.Log(secondTime)
	t.Log(milliTime)
	t.Log(nanoTime)

}

//Acquire the time according to the string
//isH true-24 hour clock false-12 hour clock
//No hours involved. The second field is ignored
func Test_GetTimeByString(t *testing.T) {
	str1 := "2020-09-16 09:21:18"
	str2 := "2020/09/16 09:21:18"
	str3 := "2020-09-16"
	str4 := "2020/09/16"
	str5 := "20200916092105"
	str6 := "20200916"

	t1_1 := GetTimeByString(str1, true)
	t1_2 := GetTimeByString(str1, false)

	t2_1 := GetTimeByString(str2, true)
	t2_2 := GetTimeByString(str2, false)

	t3_1 := GetTimeByString(str3, true) //No hours involved. The second field is ignored
	t3_2 := GetTimeByString(str3, false)

	t4_1 := GetTimeByString(str4, true) //No hours involved. The second field is ignored
	t4_2 := GetTimeByString(str4, false)

	t5_1 := GetTimeByString(str5, true)
	t5_2 := GetTimeByString(str5, false)

	t6_1 := GetTimeByString(str6, true) //No hours involved. The second field is ignored
	t6_2 := GetTimeByString(str6, false)

	t.Log(t1_1)
	t.Log(t1_2)
	t.Log(t2_1)
	t.Log(t2_2)
	t.Log(t3_1)
	t.Log(t3_2)
	t.Log(t4_1)
	t.Log(t4_2)
	t.Log(t5_1)
	t.Log(t5_2)
	t.Log(t6_1)
	t.Log(t6_2)
}

//Get the zero point of incoming time
func Test_GetZeroTime(t *testing.T) {
	zeroTime := GetZeroTime(time.Now())
	t.Log(zeroTime)
}

//Get the zero time of the day
func Test_GetNowZeroTime(t *testing.T) {
	zeroTime := GetNowZeroTime()
	t.Log(zeroTime)
}

//Get the first day of the month where the time is passed in.
// That is the zero point on the first day of a month
func Test_GetFirstDateOfMonth(t *testing.T) {
	firstDay := GetFirstDateOfMonth(time.Now())
	t.Log(firstDay)
}

//Get the last day of the month where the time is passed in
//That is, 0 o'clock on the last day of a month
func Test_GetLaDateOfMonth(t *testing.T) {
	lastDay := GetLastDateOfMonth(time.Now())
	t.Log(lastDay)
}

//Start of this year
func Test_GetFirstDateOfYear(t *testing.T) {
	first := GetFirstDateOfYear(time.Now())
	t.Log(first)
}

//End of this year
func Test_GetLastDateOfYear(t *testing.T) {
	last := GetLastDateOfYear(time.Now())
	t.Log(last)
}

//Get the day of the week for the incoming time
func Test_GetWeek(t *testing.T) {
	monday := GetWeek(GetTimeByString("2020/09/14", false))
	tuesday := GetWeek(GetTimeByString("2020/09/15", false))
	wednesday := GetWeek(GetTimeByString("2020/09/16", false))
	thursday := GetWeek(GetTimeByString("2020/09/17", false))
	friday := GetWeek(GetTimeByString("2020/09/18", false))
	saturday := GetWeek(GetTimeByString("2020/09/19", false))
	sunday := GetWeek(GetTimeByString("2020/09/20", false))

	t.Log(monday)
	t.Log(tuesday)
	t.Log(wednesday)
	t.Log(thursday)
	t.Log(friday)
	t.Log(saturday)
	t.Log(sunday)
}

//Get the month of the year
func Test_GetMonth(t *testing.T) {
	month := GetMonth(GetTimeByString("2020/09/16", false))
	t.Log(month)
}

//Is it a leap year
func Test_IsLeapYear(t *testing.T) {
	flag := IsLeapYear(GetTimeByString("2020/09/16", false))
	t.Log(flag)
}

//Time zone conversion
func Test_ParseWithLocation(t *testing.T) {
	lt, err := ParseWithLocation("America/New_York", "2020-09-16 09:21:18")
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log(lt)
}

//Convert to full time
func Test_DateTimeToString(t *testing.T) {
	dateTime := DateTimeToString(time.Now())
	t.Log(dateTime)
}

//Convert to date
func Test_DateToString(t *testing.T) {
	date := DateToString(time.Now())
	t.Log(date)
}
