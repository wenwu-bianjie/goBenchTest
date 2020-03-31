package transferTime

import (
	"fmt"
	"testing"
	"time"
)

func TestTransferTime(t *testing.T) {
	datetime := "2015-01-01 00:00:02.594581" //待转化为时间戳的字符串

	tmp := TransferTime(datetime)
	fmt.Println(tmp)
	fmt.Println(time.Now().UnixNano())
}
