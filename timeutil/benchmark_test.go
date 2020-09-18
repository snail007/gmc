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
func Benchmark_Unix(b *testing.B) {
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Unix()
	}

	b.StopTimer()
}

//Get seconds based on time
func Benchmark_TimeToUnix(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		TimeToUnix(time.Now())
	}
	b.StopTimer()
}

//Get the current time in milliseconds
func Benchmark_MilliSeconds(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		MilliSeconds()
	}
	b.StopTimer()
}

//Get the millisecond value based on time
func Benchmark_TimeToMilliSeconds(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		TimeToMilliSeconds(time.Now())
	}
	b.StopTimer()
}

//Change integer type to time type
func Benchmark_Int64ToTime(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Int64ToTime(1600307804859)
	}
	b.StopTimer()
}

//Acquire the time according to the string
//@param s timString "2006-01-02","2006/01/02 15:04:05","2006/01/02 03:04:05"
//@param isH true-24 hour clock false-12 hour clock
//No hours involved. The second field is ignored
func Benchmark_StrToTime(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		MustStrToTime("2020-09-16")
	}
	b.StopTimer()
}

//Get the zero point of incoming time
func Benchmark_ZeroTimeOf(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ZeroTimeOf(time.Now())
	}
	b.StopTimer()
}

//Get the zero time of the day
func Benchmark_ZeroTime(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ZeroTime()
	}
	b.StopTimer()
}

//Get the first day of the month where the time is passed in.
// That is the zero point on the first day of a month
func Benchmark_MonthZeroTimeOf(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		MonthZeroTimeOf(time.Now())
	}
	b.StopTimer()
}

//Get the last day of the month where the time is passed in
//That is, 0 o'clock on the last day of a month
func Benchmark_MonthLastTimeOf(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		MonthLastTimeOf(time.Now())
	}
	b.StopTimer()
}

//Start of this year
func Benchmark_YearZeroTimeOf(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		YearZeroTimeOf(time.Now())
	}
	b.StopTimer()
}

//End of this year
func Benchmark_YearLastTimeOf(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		YearLastTimeOf(time.Now())
	}
	b.StopTimer()
}

//Get the day of the week for the incoming time
func Benchmark_WeekdayOf(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		WeekdayOf(time.Now())
	}
	b.StopTimer()
}

//Get the month of the year
func Benchmark_MonthOf(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		MonthOf(time.Now())
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
func Benchmark_StrToTimeOfLocation(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		StrToTimeOfLocation("", "2020-09-17 10:01:01")
	}
	b.StopTimer()
}

//Convert to full time
func Benchmark_TimeToStr(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		TimeToStr(time.Now())
	}
	b.StopTimer()
}

//Convert to date
func Benchmark_TimeToDateStr(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		TimeToDateStr(time.Now())
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
