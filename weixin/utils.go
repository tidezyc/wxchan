package weixin

import (
	"fmt"
	"time"
)

func getR() string {
	ms := getMs()
	return fmt.Sprintf("%d", ^int32(ms))
}

func getMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
