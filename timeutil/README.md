## Demo
```golang
package main

import (
	"github.com/snail007/gmc/timeutil"
	"log"
	"time"
)

func main() {
	//Get the current time in seconds
	second := timeutil.GetNowSecond()
	log.Println("second", second)

	//Get seconds based on time
	secondByTime := timeutil.GetSecondByTime(time.Now())
	log.Println("secondByTime", secondByTime)

	//Get the current time in milliseconds
	milliSecond := timeutil.GetNowMilliSecond()
	log.Println("milliSecond", milliSecond)

	//Get the millisecond value based on time
	milliSecondByTime := timeutil.GetMilliSecondByTime(time.Now())
	log.Println("milliSecondByTime", milliSecondByTime)

	//Change integer type to time type
	timeByInt := timeutil.GetTimeByInt64(time.Now().Unix())
	log.Println("timeByInt", timeByInt)

	//Acquire the time according to the string
	//isH true-24 hour clock false-12 hour clock
	//No hours involved. The second field is ignored
	timeByString1 := timeutil.GetTimeByString("2020-09-16", false)
	timeByString2 := timeutil.GetTimeByString("2020-09-16 10:27:00", true)
	log.Println("timeByString1", timeByString1)
	log.Println("timeByString2", timeByString2)

	//Get the zero point of incoming time
	zeroTime := timeutil.GetZeroTime(time.Now())
	log.Println("zeroTime", zeroTime)

	//Get the zero time of the day
	nowZeroTime := timeutil.GetNowZeroTime()
	log.Println("nowZeroTime", nowZeroTime)

	//Get the first day of the month where the time is passed in.
	// That is the zero point on the first day of a month
	firstDateOfMonth := timeutil.GetFirstDateOfMonth(time.Now())
	log.Println("firstDateOfMonth", firstDateOfMonth)

	//Get the last day of the month where the time is passed in
	//That is, 0 o'clock on the last day of a month
	lastDateOfMonth := timeutil.GetLastDateOfMonth(time.Now())
	log.Println("lastDateOfMonth", lastDateOfMonth)

	//Start of this year
	firstDateOfYear := timeutil.GetFirstDateOfYear(time.Now())
	log.Println("firstDateOfYear", firstDateOfYear)

	//End of this year
	lastDateOfYear := timeutil.GetLastDateOfYear(time.Now())
	log.Println("lastDateOfYear", lastDateOfYear)

	//Get the day of the week for the incoming time
	week := timeutil.GetWeek(time.Now())
	log.Println("week", week)

	//Get the month of the year
	month := timeutil.GetMonth(time.Now())
	log.Println("month", month)

	//Is it a leap year
	isLeapYear := timeutil.IsLeapYear(time.Now())
	log.Println("isLeapYear", isLeapYear)


	//Time zone conversion
	parseLocate, err := timeutil.ParseWithLocation("UTC", "2020-09-17 10:23:01")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("parseLocate", parseLocate)

	//Convert to full time
	timeToString := timeutil.DateTimeToString(time.Now())
	log.Println("timeToString", timeToString)

	//Convert to date
	dateToString := timeutil.DateToString(time.Now())
	log.Println("dateToString", dateToString)

	//Time format
	timeFormat1 := timeutil.TimeFormat(time.Now(), "yyyy-dd-mm")
	timeFormat2 := timeutil.TimeFormat(time.Now(), "yyyy-dd-mm hh:ii:ss")
	timeFormat3 := timeutil.TimeFormat(time.Now(), "Y-d-m h:i:s")
	log.Println("timeFormat1", timeFormat1)
	log.Println("timeFormat2", timeFormat2)
	log.Println("timeFormat3", timeFormat3)
}

```

## Testing And Code coverage

```text
=== RUN   Test_GetNowSecond
--- PASS: Test_GetNowSecond (0.00s)
    timeutil_test.go:20: 1600310166
=== RUN   Test_GetSecondByTime
--- PASS: Test_GetSecondByTime (0.00s)
    timeutil_test.go:27: 1600310166
=== RUN   Test_GetNowMilliSecond
--- PASS: Test_GetNowMilliSecond (0.00s)
    timeutil_test.go:34: 1600310166071
=== RUN   Test_GetMilliSecondByTime
--- PASS: Test_GetMilliSecondByTime (0.00s)
    timeutil_test.go:41: 1600310166071
=== RUN   Test_GetTimeByInt64
--- PASS: Test_GetTimeByInt64 (0.04s)
    timeutil_test.go:51: 10
    timeutil_test.go:52: 13
    timeutil_test.go:53: 19
    timeutil_test.go:61: 2020-09-17 10:36:06 +0800 CST
    timeutil_test.go:62: 2020-09-17 10:36:06.071 +0800 CST
    timeutil_test.go:63: 2020-09-17 10:36:06.0717397 +0800 CST
=== RUN   Test_GetTimeByString
--- PASS: Test_GetTimeByString (0.00s)
    timeutil_test.go:109: 2020-09-16 09:21:18 +0800 CST
    timeutil_test.go:110: 2020-09-16 09:21:18 +0800 CST
    timeutil_test.go:111: 2020-09-16 09:21:18 +0800 CST
    timeutil_test.go:112: 2020-09-16 09:21:18 +0800 CST
    timeutil_test.go:113: 2020-09-16 00:00:00 +0800 CST
    timeutil_test.go:114: 2020-09-16 00:00:00 +0800 CST
    timeutil_test.go:115: 2020-09-16 00:00:00 +0800 CST
    timeutil_test.go:116: 2020-09-16 00:00:00 +0800 CST
    timeutil_test.go:117: 2020-09-16 09:21:05 +0800 CST
    timeutil_test.go:118: 2020-09-16 09:21:05 +0800 CST
    timeutil_test.go:119: 2020-09-16 00:00:00 +0800 CST
    timeutil_test.go:120: 2020-09-16 00:00:00 +0800 CST
=== RUN   Test_GetZeroTime
--- PASS: Test_GetZeroTime (0.00s)
    timeutil_test.go:159: 2020-09-17 00:00:00 +0800 CST
=== RUN   Test_GetNowZeroTime
--- PASS: Test_GetNowZeroTime (0.00s)
    timeutil_test.go:169: 2020-09-17 00:00:00 +0800 CST
=== RUN   Test_GetFirstDateOfMonth
--- PASS: Test_GetFirstDateOfMonth (0.00s)
    timeutil_test.go:180: 2020-09-01 00:00:00 +0800 CST
=== RUN   Test_GetLaDateOfMonth
--- PASS: Test_GetLaDateOfMonth (0.00s)
    timeutil_test.go:192: 2020-09-30 00:00:00 +0800 CST
=== RUN   Test_GetFirstDateOfYear
--- PASS: Test_GetFirstDateOfYear (0.00s)
    timeutil_test.go:203: 2020-01-01 00:00:00 +0800 CST
=== RUN   Test_GetLastDateOfYear
--- PASS: Test_GetLastDateOfYear (0.00s)
    timeutil_test.go:214: 2020-12-31 00:00:00 +0800 CST
=== RUN   Test_GetWeek
--- PASS: Test_GetWeek (0.00s)
    timeutil_test.go:232: 1
    timeutil_test.go:233: 2
    timeutil_test.go:234: 3
    timeutil_test.go:235: 4
    timeutil_test.go:236: 5
    timeutil_test.go:237: 6
    timeutil_test.go:238: 7
=== RUN   Test_GetMonth
--- PASS: Test_GetMonth (0.00s)
    timeutil_test.go:252: 9
=== RUN   Test_IsLeapYear
--- PASS: Test_IsLeapYear (0.00s)
    timeutil_test.go:259: true
    timeutil_test.go:263: false
=== RUN   Test_ParseWithLocation
2020/09/17 10:36:06 parsing time "panic" as "2006-01-02 15:04:05": cannot parse "panic" as "2006"
2020/09/17 10:36:06 unknown time zone panic
--- PASS: Test_ParseWithLocation (0.00s)
    timeutil_test.go:275: 2020-09-16 09:21:18 -0400 EDT
    timeutil_test.go:283: parsing time "panic" as "2006-01-02 15:04:05": cannot parse "panic" as "2006"
    timeutil_test.go:285: 0001-01-01 00:00:00 +0000 UTC
    timeutil_test.go:292: unknown time zone panic
    timeutil_test.go:294: 0001-01-01 00:00:00 +0000 UTC
=== RUN   Test_DateTimeToString
--- PASS: Test_DateTimeToString (0.00s)
    timeutil_test.go:302: 2020-09-17 10:36:06
=== RUN   Test_DateToString
--- PASS: Test_DateToString (0.00s)
    timeutil_test.go:310: 2020-09-17
=== RUN   Test_TimeFormat
--- PASS: Test_TimeFormat (0.00s)
    timeutil_test.go:328: 10:36:06
    timeutil_test.go:329: 20-09-17
    timeutil_test.go:330: 2020-09-17
    timeutil_test.go:331: 2020-09-17 10:36:06
    timeutil_test.go:332: 2020/09/1710:36:06
PASS
coverage: 96.9% of statements in ../../gmc/...

```

## Benchmark

```text
goos: windows
goarch: amd64
pkg: github.com/snail007/gmc/timeutil
Benchmark_GetNowSecond-4           	89061468	        13.7 ns/op
Benchmark_GetSecondByTime-4        	92223980	        13.4 ns/op
Benchmark_GetNowMilliSecond-4      	66702981	        16.2 ns/op
Benchmark_GetMilliSecondByTime-4   	79949898	        15.2 ns/op
Benchmark_GetTimeByInt64-4         	 2317837	       511 ns/op
Benchmark_GetTimeByString-4        	 5028524	       210 ns/op
Benchmark_GetZeroTime-4            	 3047280	       390 ns/op
Benchmark_GetNowZeroTime-4         	 2979219	       433 ns/op
Benchmark_GetFirstDateOfMonth-4    	 2219305	       536 ns/op
Benchmark_GetLastDateOfMonth-4     	 1873140	       715 ns/op
Benchmark_GetFirstDateOfYear-4     	 1808257	       581 ns/op
Benchmark_GetLastDateOfYear-4      	 1786773	       677 ns/op
Benchmark_GetWeek-4                	52200693	        22.4 ns/op
Benchmark_GetMonth-4               	25005834	        48.7 ns/op
Benchmark_IsLeapYear-4             	26670518	        45.1 ns/op
Benchmark_ParseWithLocation-4      	 3717289	       295 ns/op
Benchmark_DateTimeToString-4       	 3605575	       315 ns/op
Benchmark_DateToString-4           	 6064075	       193 ns/op
Benchmark_TimeFormat-4             	   38608	     30598 ns/op
PASS

```