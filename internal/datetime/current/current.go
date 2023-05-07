package current

import (
	"time"
)

const (
	LayoutDateTime     = "2006-01-02 15:04:05"
	LayoutDate         = "20060102"
	LayoutDateTimeNano = "2006-01-02 15:04:05.999999999"
)

type Current struct {
	timeOffSetSecond time.Duration // 偏移时间
}

func NewCurrent() *Current {
	return &Current{}
}

// SetTimeOffset 设置时间偏移
func (c *Current) SetTimeOffset(offset time.Duration) {
	c.timeOffSetSecond = offset
}

// GetTimeOffset 获取当前的时间偏移
func (c *Current) GetTimeOffset() time.Duration {
	return c.timeOffSetSecond
}

func (c *Current) Now() time.Time {
	return time.Now().Add(c.timeOffSetSecond)
}

func (c *Current) Unix() int64 {
	return time.Now().Add(c.timeOffSetSecond).Unix()
}

func (c *Current) UnixNano() int64 {
	return time.Now().Add(c.timeOffSetSecond).UnixNano()
}

func (c *Current) UnixMill() int64 {
	return time.Now().Add(c.timeOffSetSecond).UnixMilli()
}
