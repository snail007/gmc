/*
 * @Description:Time tools test
 * @Author: bearcat-panda
 * @Date: 2020-09-16 8:00:24
 * @LastEditTime: 2020-09-16 10:36:24
 * @LastEditors: bearcat-panda
 */
package timeutil

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

//Get the current time in seconds
func Test_GetNowSecond(t *testing.T) {
	nowSecond := GetNowSecond()
	t.Log(nowSecond)
	assert.Len(t, fmt.Sprintf("%d", nowSecond), 10)
}

//Get seconds based on time
func Test_GetSecondByTime(t *testing.T) {
	second := GetSecondByTime(time.Now())
	t.Log(second)
	assert.Len(t, fmt.Sprintf("%d", second), 10)
}

//Get the current time in milliseconds
func Test_GetNowMilliSecond(t *testing.T) {
	nowMilliSecond := GetNowMilliSecond()
	t.Log(nowMilliSecond)
	assert.Len(t, fmt.Sprintf("%d", nowMilliSecond), 13)
}

//Get the millisecond value based on time
func Test_GetMilliSecondByTime(t *testing.T) {
	milliSecond := GetMilliSecondByTime(time.Now())
	t.Log(milliSecond)
	assert.Len(t, fmt.Sprintf("%d", milliSecond), 13)
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

	assert.NotNil(t, secondTime)
	assert.Equal(t, secondTime.Format("2006-01-02 15:04:05"),
		time.Now().Format("2006-01-02 15:04:05"))
	assert.NotNil(t, milliTime)
	assert.Equal(t, milliTime.Format("2006-01-02 15:04:05"),
		time.Now().Format("2006-01-02 15:04:05"))
	assert.NotNil(t, nanoTime)
	assert.Equal(t, nanoTime.Format("2006-01-02 15:04:05"),
		time.Now().Format("2006-01-02 15:04:05"))

	assert.Panics(t, func() { //panic
		GetTimeByInt64(123)
	}, "Length does not meet the requirements")
}

//Acquire the time according to the string
//@param s timString "2006-01-02","2006/01/02 15:04:05","2006/01/02 03:04:05"
//@param isH true-24 hour clock false-12 hour clock
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

	assert.NotNil(t, t1_1)
	assert.Len(t, t1_1.Format("2006-01-02 15:04:05"), len(str1))
	assert.NotNil(t, t1_2)
	assert.Len(t, t1_2.Format("2006-01-02 15:04:05"), len(str1))
	assert.NotNil(t, t2_1)
	assert.Len(t, t2_1.Format("2006-01-02 15:04:05"), len(str2))
	assert.NotNil(t, t2_2)
	assert.Len(t, t2_2.Format("2006-01-02 15:04:05"), len(str2))
	assert.NotNil(t, t3_1)
	assert.Len(t, t3_1.Format("2006-01-02"), len(str3))
	assert.NotNil(t, t3_2)
	assert.Len(t, t3_2.Format("2006-01-02"), len(str3))
	assert.NotNil(t, t4_1)
	assert.Len(t, t4_1.Format("2006-01-02"), len(str4))
	assert.NotNil(t, t4_2)
	assert.Len(t, t4_2.Format("2006-01-02"), len(str4))
	assert.NotNil(t, t5_1)
	assert.Len(t, t5_1.Format("20060102150405"), len(str5))
	assert.NotNil(t, t5_2)
	assert.Len(t, t5_2.Format("20060102150405"), len(str5))
	assert.NotNil(t, t6_1)
	assert.Len(t, t6_1.Format("20060102"), len(str6))
	assert.NotNil(t, t6_2)
	assert.Len(t, t6_2.Format("20060102"), len(str6))

	assert.Panics(t, func() { //panic
		GetTimeByString("", true)
	}, "String format error")

	assert.Panics(t, func() { //panic
		GetTimeByString("1", true)
	}, "String format error")
}

//Get the zero point of incoming time
func Test_GetZeroTime(t *testing.T) {
	zeroTime := GetZeroTime(time.Now())
	t.Log(zeroTime)

	assert.NotNil(t, zeroTime)
	assert.Equal(t, zeroTime.Format("2006-01-02"),
		time.Now().Format("2006-01-02"))
}

//Get the zero time of the day
func Test_GetNowZeroTime(t *testing.T) {
	zeroTime := GetNowZeroTime()
	t.Log(zeroTime)

	assert.NotNil(t, zeroTime)
	assert.Equal(t, zeroTime.Format("2006-01-02"),
		time.Now().Format("2006-01-02"))
}

//Get the first day of the month where the time is passed in.
// That is the zero point on the first day of a month
func Test_GetFirstDateOfMonth(t *testing.T) {
	firstDay := GetFirstDateOfMonth(time.Now())
	t.Log(firstDay)
	assert.NotNil(t, firstDay)

	now := time.Now()
	assert.Equal(t, firstDay.Format("2006-01-02"),
		now.AddDate(0, 0, -now.Day()+1).Format("2006-01-02"))
}

//Get the last day of the month where the time is passed in
//That is, 0 o'clock on the last day of a month
func Test_GetLaDateOfMonth(t *testing.T) {
	lastDay := GetLastDateOfMonth(time.Now())
	t.Log(lastDay)
	assert.NotNil(t, lastDay)

	now := time.Now()
	assert.Equal(t, lastDay.Format("2006-01-02"),
		now.AddDate(0, 1, -now.Day()).Format("2006-01-02"))
}

//Start of this year
func Test_GetFirstDateOfYear(t *testing.T) {
	first := GetFirstDateOfYear(time.Now())
	t.Log(first)
	assert.NotNil(t, first.String())

	now := time.Now()
	assert.Equal(t, first.Format("2006-01-02"),
		now.AddDate(0, -int(now.Month())+1, -now.Day()+1).Format("2006-01-02"))
}

//End of this year
func Test_GetLastDateOfYear(t *testing.T) {
	last := GetLastDateOfYear(time.Now())
	t.Log(last)

	assert.NotNil(t, last)
	now := time.Now()
	assert.Equal(t, last.Format("2006-01-02"),
		now.AddDate(1, -int(now.Month())+1, -now.Day()).Format("2006-01-02"))
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

	assert.Equal(t, monday, 1)
	assert.Equal(t, tuesday, 2)
	assert.Equal(t, wednesday, 3)
	assert.Equal(t, thursday, 4)
	assert.Equal(t, friday, 5)
	assert.Equal(t, saturday, 6)
	assert.Equal(t, sunday, 7)
}

//Get the month of the year
func Test_GetMonth(t *testing.T) {
	month := GetMonth(GetTimeByString("2020/09/16", false))
	t.Log(month)
	assert.Equal(t, month, 9)
}

//Is it a leap year
func Test_IsLeapYear(t *testing.T) {
	flag := IsLeapYear(GetTimeByString("2020/09/16", false))
	t.Log(flag)
	assert.Equal(t, flag, true)

	flag = IsLeapYear(GetTimeByString("2007/09/16", false))
	t.Log(flag)
	assert.Equal(t, flag, false)
}

//Time zone conversion
//@param location "","UTC","merica/New_York"
//@param timeStrEx	"2006-01-02 15:04:05"
func Test_ParseWithLocation(t *testing.T) {
	lt, err := ParseWithLocation("America/New_York", "2020-09-16 09:21:18")
	if err != nil {
		t.Log(err)
		return
	}

	t.Log(lt)

	assert.NotNil(t, lt)
	assert.Len(t, lt.Format("2006-01-02 15:04:05"), len("2020-09-16 09:21:18"))

	//panic
	lt, err = ParseWithLocation("", "panic")
	if err != nil {
		t.Log(err)
	}
	t.Log(lt)

	assert.Error(t, err, "unknown time zone panic")

	//panic
	lt, err = ParseWithLocation("panic", "panic")
	if err != nil {
		t.Log(err)
	}
	t.Log(lt)

	assert.Error(t, err, "unknown time zone panic")
}

//Convert to full time
func Test_DateTimeToString(t *testing.T) {
	dateTime := DateTimeToString(time.Now())
	t.Log(dateTime)

	assert.Len(t, dateTime, 19)
}

//Convert to date
func Test_DateToString(t *testing.T) {
	date := DateToString(time.Now())
	t.Log(date)
	assert.Len(t, date, 10)
}

//Time format
//@param format  "h:i:s","Y-m-d h:i:s","y-m-d"
func Test_TimeFormat(t *testing.T) {
	format1 := "h:i:s"
	format2 := "y-m-d"
	format3 := "Y-m-d"
	format4 := "Y-m-d h:i:s"
	format5 := "Y/m/dh:i:s"

	s1 := TimeFormat(time.Now(), format1)
	s2 := TimeFormat(time.Now(), format2)
	s3 := TimeFormat(time.Now(), format3)
	s4 := TimeFormat(time.Now(), format4)
	s5 := TimeFormat(time.Now(), format5)

	t.Log(s1)
	t.Log(s2)
	t.Log(s3)
	t.Log(s4)
	t.Log(s5)

	assert.Len(t, s1, 8)
	assert.Len(t, s2, 8)
	assert.Len(t, s3, 10)
	assert.Len(t, s4, 19)
	assert.Len(t, s5, 18)
}
