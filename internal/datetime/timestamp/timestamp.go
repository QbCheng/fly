package timestamp

import (
	"fly/internal/datetime/timeUtil"
	"strconv"
	"time"
)

const (
	UnitSecondsToMinute = 60
	UnitMinutesToHour   = 60
	UnitHoursToDay      = 24

	UnitSecondsToHour = UnitSecondsToMinute * UnitMinutesToHour

	UnitSecondsToDay = UnitSecondsToHour * UnitHoursToDay
	UnitMinutesToDay = UnitMinutesToHour * UnitHoursToDay

	UnitMsToMinute = UnitSecondsToMinute * 1000
	UnitMsToHour   = UnitSecondsToHour * 1000
	UnitMsToDay    = UnitSecondsToDay * 1000
	UnitMsToSecond = 1000
)

// ParseTime 时间戳 转化 time.Timer
func ParseTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

// IsSameMinute 同一分钟
func IsSameMinute(t1, t2 int64) bool {
	return t1/UnitSecondsToMinute == t2/UnitSecondsToMinute
}

// IsSameHour 同一小时
func IsSameHour(t1, t2 int64) bool {
	return t1/UnitSecondsToHour == t2/UnitSecondsToHour
}

// DifferenceDays 相差的天数
func DifferenceDays(t1, t2 int64) int {
	if t1 > t2 {
		t1, t2 = t2, t1
	}

	d := (t2 - t1) / UnitSecondsToDay
	t := t1 + UnitSecondsToDay*d
	if !IsSameDay(t, t2) {
		d++
	}

	return int(d)
}

// IsSameDay 同一天检测
func IsSameDay(t1, t2 int64) bool {
	time1 := time.Unix(t1, 0)
	time2 := time.Unix(t2, 0)
	return time1.YearDay() == time2.YearDay() && time1.Year() == time2.Year()
}

// IsSameWeek 是同一周
func IsSameWeek(t1, t2 int64) bool {
	lt1 := t1
	lt2 := t2
	tt1 := time.Unix(lt1, 0)
	tt2 := time.Unix(lt2, 0)
	y1, w1 := tt1.ISOWeek()
	y2, w2 := tt2.ISOWeek()
	return y1 == y2 && w1 == w2
}

// IsSameMonth 是同一月
func IsSameMonth(t1, t2 int64) bool {
	tt1 := time.Unix(t1, 0)
	tt2 := time.Unix(t2, 0)
	return tt1.Year() == tt2.Year() && tt1.Month() == tt2.Month()
}

// CurrentDayZero 获得当前时间的零点, 通过时间戳
func CurrentDayZero(timestamp int64) int64 {
	return timeUtil.CurrentDayZero(ParseTime(timestamp)).Unix()
}

// MonthStart 月份开始
func MonthStart(timestamp int64) int64 {
	return timeUtil.MonthStart(ParseTime(timestamp)).Unix()
}

// YearStart 通过当年1月1日0点时间戳
func YearStart(timestamp int64) int64 {
	return timeUtil.YearStart(ParseTime(timestamp)).Unix()
}

// Weekday 返回本周的周几
func Weekday(timestamp int64) time.Weekday {
	return ParseTime(timestamp).Weekday()
}

// TransformWeekday 将 数值 转化为 time.Weekday
func TransformWeekday(weekday int64) time.Weekday {
	return time.Weekday(weekday % 7)
}

// NowWeekdayWithStart 本周, 指定星期的开始时间
func NowWeekdayWithStart(timestamp int64, weekday time.Weekday) int64 {
	s := CurrentDayZero(timestamp)
	currentWeekday := Weekday(timestamp)
	if weekday == currentWeekday {
		return s
	}
	if weekday == time.Sunday {
		weekday = 7
	}
	if currentWeekday == time.Sunday {
		currentWeekday = 7
	}
	return s + int64(weekday-currentWeekday)*86400
}

// SpecifiedMonthStartAndEnd 获得 某年某月 开始时间 和 结束时间
func SpecifiedMonthStartAndEnd(y int, m int) (int64, int64) {
	// 数字月份必须前置补零
	var ms string
	if m < 10 {
		ms = "0" + strconv.Itoa(m)
	} else {
		ms = strconv.Itoa(m)
	}
	theTime, _ := time.ParseInLocation("2006-01-02 15:04:05", strconv.Itoa(y)+"-"+ms+"-01 00:00:00", time.Local)
	t1 := time.Date(y, theTime.Month(), 1, 0, 0, 0, 0, time.Local)
	t2 := time.Date(y, theTime.Month()+1, 1, 0, 0, 0, 0, time.Local).Add(-time.Second)
	return t1.Unix(), t2.Unix()
}
