package current

import (
	"time"
)

var static = NewCurrent()

// SetTimeOffset 设置时间偏移
func SetTimeOffset(offset time.Duration) {
	static.SetTimeOffset(offset)
}

// GetTimeOffset 获取当前的时间偏移
func GetTimeOffset() time.Duration {
	return static.GetTimeOffset()
}

func Now() time.Time {
	return time.Now().Add(static.timeOffSetSecond)
}

func Unix() int64 {
	return time.Now().Add(static.timeOffSetSecond).Unix()
}

func UnixNano() int64 {
	return time.Now().Add(static.timeOffSetSecond).UnixNano()
}

func UnixMill() int64 {
	return time.Now().Add(static.timeOffSetSecond).UnixMilli()
}
