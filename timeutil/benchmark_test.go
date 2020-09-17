/*
 * @Description:Time tools
 * @Author: bearcat-panda
 * @Date: 2020-09-16 8:00:24
 * @LastEditTime: 2020-09-17 10:36:24
 * @LastEditors: bearcat-panda
 */

package timeutil

import (
	"testing"
	"time"
)

//Get the current time in seconds
func Benchmark_GetNowSecond(b *testing.B) {
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GetNowSecond()
	}

	b.StopTimer()
}

//Get seconds based on time
func Benchmark_GetSecondByTime(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GetSecondByTime(time.Now())
	}
	b.StopTimer()
}

//Get the current time in milliseconds
func Benchmark_GetNowMilliSecond(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GetNowMilliSecond()
	}
	b.StopTimer()
}

//Get the millisecond value based on time
func Benchmark_GetMilliSecondByTime(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GetMilliSecondByTime(time.Now())
	}
	b.StopTimer()
}

//Change integer type to time type
func Benchmark_GetTimeByInt64(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GetTimeByInt64(1600307804859)
	}
	b.StopTimer()
}

//Acquire the time according to the string
//@param s timString "2006-01-02","2006/01/02 15:04:05","2006/01/02 03:04:05"
//@param isH true-24 hour clock false-12 hour clock
//No hours involved. The second field is ignored
func Benchmark_GetTimeByString(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GetTimeByString("2020-09-16", false)
	}
	b.StopTimer()
}

//Get the zero point of incoming time
func Benchmark_GetZeroTime(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GetZeroTime(time.Now())
	}
	b.StopTimer()
}

//Get the zero time of the day
func Benchmark_GetNowZeroTime(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GetNowZeroTime()
	}
	b.StopTimer()
}

//Get the first day of the month where the time is passed in.
// That is the zero point on the first day of a month
func Benchmark_GetFirstDateOfMonth(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GetFirstDateOfMonth(time.Now())
	}
	b.StopTimer()
}

//Get the last day of the month where the time is passed in
//That is, 0 o'clock on the last day of a month
func Benchmark_GetLastDateOfMonth(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GetLastDateOfMonth(time.Now())
	}
	b.StopTimer()
}

//Start of this year
func Benchmark_GetFirstDateOfYear(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GetFirstDateOfYear(time.Now())
	}
	b.StopTimer()
}

//End of this year
func Benchmark_GetLastDateOfYear(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GetLastDateOfYear(time.Now())
	}
	b.StopTimer()
}

//Get the day of the week for the incoming time
func Benchmark_GetWeek(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GetWeek(time.Now())
	}
	b.StopTimer()
}

//Get the month of the year
func Benchmark_GetMonth(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GetMonth(time.Now())
	}
	b.StopTimer()
}

//Is it a leap year
func Benchmark_IsLeapYear(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		IsLeapYear(time.Now())
	}
	b.StopTimer()
}

//Time zone conversion
//@param location "","UTC","merica/New_York"
//@param timeStrEx	"2006-01-02 15:04:05"
func Benchmark_ParseWithLocation(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ParseWithLocation("", "2020-09-17 10:01:01")
	}
	b.StopTimer()
}

//Convert to full time
func Benchmark_DateTimeToString(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		DateTimeToString(time.Now())
	}
	b.StopTimer()
}

//Convert to date
func Benchmark_DateToString(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		DateToString(time.Now())
	}
	b.StopTimer()
}

//Time format
//@param format  "h:i:s","Y-m-d h:i:s","y-m-d"
func Benchmark_TimeFormat(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		TimeFormat(time.Now(), "2006-09-16")
	}
	b.StopTimer()
}
