package epss

import (
	"fmt"
	"time"
)

// db key 的分隔符
const dbKeySep = "."

// 生成 用户名 检索 前缀
func makeDbMsgPre(to string) []byte {
	pre := make([]byte, 4+len(to)+1)
	bp := copy(pre, "msg")
	bp += copy(pre[bp:], dbKeySep)
	bp += copy(pre[bp:], to)
	bp += copy(pre[bp:], dbKeySep)
	return pre
}

// makeDbMsgKey 生成 用户名 消息 存储 key
// 由 user + 时间 排序，from+id 定位
func makeDbMsgKey(to string, t *time.Time, src uint8, id string) []byte {
	srcStr := fmt.Sprintf("%02X", src)
	//fromStr := strconv.FormatUint(uint64(from), 16)
	if t == nil {
		now := time.Now()
		t = &now
	}
	tStr := timeToStr(t)
	key := make([]byte, 4+len(to)+len(tStr)+len(srcStr)+len(id)+3)
	bp := copy(key, "msg") // msg 前缀
	bp += copy(key[bp:], dbKeySep)
	bp += copy(key[bp:], to)
	bp += copy(key[bp:], dbKeySep)
	bp += copy(key[bp:], tStr)
	bp += copy(key[bp:], dbKeySep)
	bp += copy(key[bp:], srcStr)
	bp += copy(key[bp:], dbKeySep)
	bp += copy(key[bp:], id)
	return key
}

func timeToStr(t *time.Time) string {
	return fmt.Sprintf("%04d%02d%02d%02d%02d%02d%03d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1000000)
}
