package epss

import "strconv"

// 用户名 + 客户端kind 唯一标示
func makeCliKey(name string, kind uint8) string {
	return name + strconv.Itoa(int(kind))
}

// 生成 用户名 检索 前缀
func makeDbPre(to string, cli uint8) []byte {
	kind := strconv.Itoa(int(cli))
	pre := make([]byte, len(to)+len(kind)+1)
	bp := copy(pre, to)
	bp += copy(pre[bp:], "|")
	bp += copy(pre[bp:], kind)
	// bp += copy(pre[bp:], "|")
	return pre
}

// 生成 用户名 消息 存储 key
func makeDbKey(to string, cli uint8, id string) []byte {
	kind := strconv.Itoa(int(cli))
	key := make([]byte, len(to)+len(kind)+len(id)+2)
	bp := copy(key, to)
	bp += copy(key[bp:], "|")
	bp += copy(key[bp:], kind)
	bp += copy(key[bp:], "|")
	bp += copy(key[bp:], id)
	return key
}
