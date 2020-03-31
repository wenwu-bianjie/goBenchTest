package transferTime

import (
	"strconv"
	"strings"
	"time"
)

func TransferTime(time_str string) int64 {
	loc, _ := time.LoadLocation("Local") //获取时区
	tmp, _ := time.ParseInLocation("2006-01-02 15:04:05", time_str, loc)
	timestamp := tmp.Unix() * 1000000000 //转化为时间戳 类型是int64
	millisecondArr := strings.SplitN(time_str, ".", -1)
	if len(millisecondArr) > 1 {
		if millisecond, err := strconv.Atoi(millisecondArr[1]); err == nil {
			return timestamp + (int64(millisecond) * 1000)
		}
	}
	return timestamp
}
