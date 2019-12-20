package epss

import (
	"container/list"
	"epssd/json"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"time"

	"github.com/gorilla/websocket"
)

// Pusher 往 客户端 推送 接口
type Pusher interface {
	pushString(msg string)
	pushBytes(msg []byte)
	logout(tip string)
}

// 客户端 类型

const (
	// KindApp 客户端 app 主 消息 提示程序
	KindApp uint8 = 1
	// KindEmail web应用 邮件
	KindEmail uint8 = 2
	// KindDocsy web应用 文档
	KindDocsy uint8 = 3
)

type msgInfo struct {
	logout bool
	msg    []byte
}

// CliConn 客户连接
type CliConn struct {
	s    *Server
	kind uint8  // 客户端类型
	name string // 客户名称

	conn   *websocket.Conn
	ch     chan *msgInfo
	closed bool // 用来 标记
}

func newCliConn(s *Server, conn *websocket.Conn) *CliConn {
	c := &CliConn{
		s:      s,
		conn:   conn,
		ch:     make(chan *msgInfo, 10),
		closed: false,
	}
	return c
}

func (c *CliConn) pushString(msg string) {
	info := &msgInfo{
		logout: false,
		msg:    []byte(msg),
	}
	c.ch <- info
}

func (c *CliConn) pushBytes(msg []byte) {
	info := &msgInfo{
		logout: false,
		msg:    make([]byte, len(msg)),
	}
	copy(info.msg, msg)
	c.ch <- info
}

func (c *CliConn) logout(tip string) {
	info := &msgInfo{
		logout: true,
	}
	info.msg = websocket.FormatCloseMessage(websocket.CloseNormalClosure, tip)
	c.ch <- info
}

// handleSend 长连接 里的 待发送消息 队列
func (c *CliConn) handle() {
	for {
		msg := <-c.ch
		if msg == nil {
			// 发送空 来 判断 是否 退出
			break
		}

		if msg.logout {
			// 主动 关闭 然后 read 里 就会 报错
			_ = c.conn.WriteControl(websocket.CloseMessage, msg.msg, time.Time{})
			//break
			// 这里 不推出 还是由 read 里 报错 在 <- nil 退出
			continue
		}
		err := c.conn.WriteMessage(websocket.TextMessage, msg.msg)
		if err != nil {
			c.s.log.Error("handleSend.Write err:", err)
		}
	}
	close(c.ch)
	c.closed = true
	c.s.log.Debug("the Send close:", c.conn.RemoteAddr().String())
}

// Serve 管理 客户端接入 的 长连接
func (c *CliConn) Serve() {
	c.s.log.Debug("the websocket connect:", c.conn.RemoteAddr().String())

	// 新建一个 gorou 用来 并发 发送 消息
	go c.handle()

	for {
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			if _, ok := err.(*websocket.CloseError); !ok {
				c.s.log.Error("ReadMessage err:", err)
			}
			c.ch <- nil // 发送 空 使 gorou 正常 退出
			break
		}
		// 解析 客户端 连接 登录
		j, err := json.LoadBytes(p)
		if err != nil {
			c.s.log.Error("异常连接, not json:", string(p))
			c.logout("notjson")
			continue
		}

		fun := j.GetString("fun")
		if len(fun) == 0 {
			c.s.log.Error("异常连接, err json", j)
			c.logout("errjson")
			continue
		}
		switch fun {
		case "login":
			c.login(j)
		default:
			c.logout("errfun")
		}
	}
	_ = c.conn.Close()
	c.s.log.Debug("the websocket disconn:", c.conn.RemoteAddr().String())
}

func (c *CliConn) login(j *json.JSON) {
	name := j.GetString("name")
	cli := j.GetUint8("cli")

	// 添加 验证 用户

	// 添加 返回信息
	c.pushString(`{"ret": "login", "text": "success"}`)
	if len(c.name) > 0 { // 已经登录过了
		return
	}

	c.s.setAndKickUser(name, cli, c)

	// 读取 db 里 未读 的 消息
	delList := list.New()
	preKey := makeDbPre(name, cli)
	ite := c.s.db.NewIterator(util.BytesPrefix(preKey), nil)
	for ite.Next() {
		c.pushBytes(ite.Value())
		buf := make([]byte, len(ite.Key()))
		copy(buf, ite.Key())
		delList.PushBack(buf)
	}
	ite.Release()

	// 批量 删除
	batch := new(leveldb.Batch)
	for i := delList.Front(); i != nil; i = i.Next() {
		if k, o := i.Value.([]byte); o {
			batch.Delete(k)
		}
	}
	err := c.s.db.Write(batch, nil)
	if err != nil {
		c.s.log.Error("Write batch err:", err)
	}
}
