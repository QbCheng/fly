package timestamp

import (
	"testing"
	"time"
)

func TestIsSameMinute(t *testing.T) {
	type args struct {
		t1 int64
		t2 int64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "2023-03-28 19:41:00 与 2023-03-29 19:41:00 不是同一分钟",
			args: args{
				t1: 1680003660,
				t2: 1680003660 + 86400,
			},
			want: false,
		},
		{
			name: "2023-03-28 19:41:00 与 2023-03-28 19:42:00 不是同一分钟",
			args: args{
				t1: 1680003660,
				t2: 1680003720,
			},
			want: false,
		},
		{
			name: "2023-03-28 19:41:00 与 2023-03-28 19:41:59 是同一分钟",
			args: args{
				t1: 1680003660,
				t2: 1680003719,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSameMinute(tt.args.t1, tt.args.t2); got != tt.want {
				t.Errorf("IsSameMinuteByTimestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSameHour(t *testing.T) {
	type args struct {
		t1 int64
		t2 int64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "2023-03-28 19:00:00 与 2023-03-28 19:00:00 是同一小时",
			args: args{
				t1: 1680001200,
				t2: 1680001200 + 3599,
			},
			want: true,
		},
		{
			name: "2023-03-28 19:00:00 与 2023-03-28 20:00:00 不是同一小时",
			args: args{
				t1: 1680001200,
				t2: 1680001200 + 3600,
			},
			want: false,
		},
		{
			name: "2023-03-28 19:00:00 与 2024-03-28 19:00:00 不是同一小时",
			args: args{
				t1: 1680001200,
				t2: 1680001200 + 86400,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSameHour(tt.args.t1, tt.args.t2); got != tt.want {
				t.Errorf("IsSameMinuteByTimestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCurrentDayZero(t *testing.T) {
	type args struct {
		timestamp int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "2023-03-28 19:00:00 当前天的第一秒",
			args: args{
				timestamp: 1680001200,
			},
			want: 1679932800,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CurrentDayZero(tt.args.timestamp); got != tt.want {
				t.Errorf("CurrentDayZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDifferenceDays(t *testing.T) {
	type args struct {
		t1 int64
		t2 int64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "2023-03-28 19:00:00 和 2023-04-28 19:00:00 相差 31天",
			args: args{
				t1: 1680001200,
				t2: 1682611200,
			},
			want: 31,
		},
		{
			name: "2023-03-28 19:00:00 和 2023-03-28 20:00:00 相差 0天",
			args: args{
				t1: 1680001200,
				t2: 1680001200 + 3600,
			},
			want: 0,
		},
		{
			name: "2023-03-28 19:00:00 和 2023-03-29 00:00:00 相差 1天",
			args: args{
				t1: 1680001200,
				t2: 1680001200 + 3600*5,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DifferenceDays(tt.args.t1, tt.args.t2); got != tt.want {
				t.Errorf("DifferenceDays() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSameWeek(t *testing.T) {
	type args struct {
		t1 int64
		t2 int64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "2023-03-28 19:00:00 和 2023-03-26 19:00:00 不是同一周",
			args: args{
				t1: 1680001200,
				t2: 1680001200 - 86400*2,
			},
			want: false,
		},
		{
			name: "2023-03-28 19:00:00 和 2023-03-27 19:00:00 是同一周",
			args: args{
				t1: 1680001200,
				t2: 1680001200 - 86400,
			},
			want: true,
		},
		{
			name: "2023-03-28 19:00:00 和 2023-04-02 19:00:00 是同一周",
			args: args{
				t1: 1680001200,
				t2: 1680001200 + 86400*4,
			},
			want: true,
		},
		{
			name: "2023-03-28 19:00:00 和 2023-04-03 19:00:00 是同一周",
			args: args{
				t1: 1680001200,
				t2: 1680001200 + 86400*5,
			},
			want: true,
		},
		{
			name: "2023-03-28 19:00:00 和 2023-04-04 19:00:00 不是同一周",
			args: args{
				t1: 1680001200,
				t2: 1680001200 + 86400*6,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSameWeek(tt.args.t1, tt.args.t2); got != tt.want {
				t.Errorf("IsSameWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMonthStart(t *testing.T) {
	type args struct {
		timestamp int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "2023-03-28 19:00:00",
			args: args{
				timestamp: 1680001200,
			},
			want: 1677600000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MonthStart(tt.args.timestamp); got != tt.want {
				t.Errorf("MonthStart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNowWeekdayWithStart(t *testing.T) {
	type args struct {
		timestamp int64
		weekday   time.Weekday
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "2023-03-28 19:00:00 时间的周六",
			args: args{
				timestamp: 1680001200,
				weekday:   time.Saturday,
			},
			want: 1680278400,
		},
		{
			name: "2023-03-28 19:00:00 时间的周一",
			args: args{
				timestamp: 1680001200,
				weekday:   time.Monday,
			},
			want: 1679846400,
		},
		{
			name: "2023-03-28 19:00:00 时间的周日",
			args: args{
				timestamp: 1680001200,
				weekday:   time.Sunday,
			},
			want: 1680364800,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NowWeekdayWithStart(tt.args.timestamp, tt.args.weekday); got != tt.want {
				t.Errorf("NowWeekdayWithStart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReverseTransformWeekday(t *testing.T) {
	type args struct {
		weekday int64
	}
	tests := []struct {
		name string
		args args
		want time.Weekday
	}{
		{
			name: "数值1为周一",
			args: args{
				weekday: 1,
			},
			want: time.Monday,
		},
		{
			name: "数值0为周日",
			args: args{
				weekday: 0,
			},
			want: time.Sunday,
		},
		{
			name: "数值7为周日",
			args: args{
				weekday: 7,
			},
			want: time.Sunday,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TransformWeekday(tt.args.weekday); got != tt.want {
				t.Errorf("ReverseTransformWeekday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWeekday(t *testing.T) {
	type args struct {
		timestamp int64
	}
	tests := []struct {
		name string
		args args
		want time.Weekday
	}{
		{
			name: "2023-03-28 19:00:00 时间",
			args: args{
				timestamp: 1680001200,
			},
			want: time.Tuesday,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Weekday(tt.args.timestamp); got != tt.want {
				t.Errorf("Weekday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYearStart(t *testing.T) {
	type args struct {
		timestamp int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "2023-03-28 19:00:00 时间",
			args: args{
				timestamp: 1680001200,
			},
			want: 1672502400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := YearStart(tt.args.timestamp); got != tt.want {
				t.Errorf("YearStart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpecifiedMonthStartAndEnd(t *testing.T) {
	type args struct {
		y int
		m int
	}
	tests := []struct {
		name  string
		args  args
		want  int64
		want1 int64
	}{
		{
			name: "2023年3月 的第一秒和最后一秒",
			args: args{
				y: 2023,
				m: 3,
			},
			want:  1677600000,
			want1: 1680278399,
		},
		{
			name: "2023年4月 的第一秒和最后一秒",
			args: args{
				y: 2023,
				m: 4,
			},
			want:  1680278400,
			want1: 1682870399,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := SpecifiedMonthStartAndEnd(tt.args.y, tt.args.m)
			if got != tt.want {
				t.Errorf("SpecifiedMonthStartAndEnd() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SpecifiedMonthStartAndEnd() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
