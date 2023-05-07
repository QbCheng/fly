package timeUtil

import "time"

// DifferenceDaysByTimer 相差天数
func DifferenceDaysByTimer(t1, t2 time.Time) int {
	if t1.After(t2) {
		t1, t2 = t2, t1
	}
	// 计算 相隔天数
	d := t2.Sub(t1) / (time.Hour * 24)

	t := t1.Add(d * time.Hour * 24)
	if !IsSameDay(t, t2) {
		d++
	}
	return int(d)
}

// IsSameDay 同一天检测
func IsSameDay(t1, t2 time.Time) bool {
	return t1.YearDay() == t2.YearDay() && t1.Year() == t2.Year()
}

// IsSameWeek 同一周
func IsSameWeek(t1, t2 time.Time) bool {
	y1, w1 := t1.ISOWeek()
	y2, w2 := t2.ISOWeek()
	return y1 == y2 && w1 == w2
}

// IsSameMonth 同一个月
func IsSameMonth(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month()
}

// InTimeRange 在时间范围内
func InTimeRange(now, bT, eT time.Time) bool {
	return (now.Equal(bT) || now.After(bT)) && now.Before(eT)
}

// CurrentDayZero 获得当前时间的零点
func CurrentDayZero(now time.Time) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

// MonthStart 获得当前时间的月份的第一天
func MonthStart(now time.Time) time.Time {
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
}

// YearStart 获得当前时间的年份的第一天
func YearStart(now time.Time) time.Time {
	return time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, now.Location())
}
