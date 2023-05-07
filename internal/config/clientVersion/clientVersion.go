package clientVersion

import (
	"errors"
	"strconv"
	"strings"
)

const (
	//	1.0.1 形式的版本号
	MainVersion    = 0x3
	MinorVersion   = 0x3
	RevisedVersion = 0xFFF
	SegmentNumber  = 3
)

var segmentSize = []uint16{RevisedVersion, MinorVersion, MainVersion}
var segmentMove = []uint16{12, 2, 2}

// ParseClientVersion 解析客户端版本
func ParseClientVersion(id uint16) (string, error) {
	ret := ""
	for i := 0; i < 3; i++ {
		tempInt := id & segmentSize[i]
		ret = strconv.Itoa(int(tempInt)) + "." + ret
		id = id >> segmentMove[i]
	}
	ret = strings.TrimSuffix(ret, ".")
	return ret, nil
}

func FormatClientVersion(id string) (uint16, error) {
	// 分割
	segments := strings.Split(id, ".")
	if len(segments) != SegmentNumber {
		return 0, errors.New("Invalid client version. ")
	}
	segments = segments[:SegmentNumber]
	var ret uint16 = 0
	var pos uint16 = 16
	for i, segment := range segments {
		pos -= segmentMove[len(segmentMove)-i-1]
		tempInt, err := strconv.Atoi(segment)
		if err != nil {
			return 0, err
		}
		if tempInt < 0 {
			return 0, errors.New("Invalid client version. ")
		}
		tempInt = tempInt << pos
		ret |= uint16(tempInt)
	}
	return ret, nil
}
