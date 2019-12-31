package epss

import (
	"epssd/json"
	"github.com/syndtr/goleveldb/leveldb/util"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

// pusher 往 客户端 推送 接口
type pusher interface {
	pushString(msg string) bool
	// variable true 表示 msg []byte 内容会变，需要做copy处理 false 就不需要
	pushBytes(msg []byte, variable bool) bool
	logout(tip string)
}

const (
	// 客户端 类型

	// KindApp 客户端 app 主 消息 提示程序
	KindApp uint8 = 1 // 1是 主客户端 可以获取 所有类型客户端 的 消息

	// 页面客户端 只能接收 该页面 发送的消息

	// KindEmail web应用 邮件
	KindEmail uint8 = 0xA
	// KindDocsy web应用 文档
	KindDocsy uint8 = 0xB

	// ox8X 给第三方 集成使用

	// 客户端 Flags 说明

	// FlagMessage 正常的消息流转
	FlagMessage uint32 = 0x0001 // 消息
	// FlagNotice 通知 布告，对应所有人
	FlagNotice uint32 = 0x0002 // 通知 告示 布告 公告
	// FlagBroad 广播 及时性 的 通知 不会存储db
	FlagBroad uint32 = 0x0004 // 广播

)

// the ConnState type are defined the CliConn state
const (
	NullConnState    int32 = 0
	RunningConnState int32 = 1
	ClosingConnState int32 = 2
	ClosedConnState  int32 = 3
)

// the send message type are defined sendInfo.t
const (
	// LogoutMsg 注销 退出
	LogoutMsg uint8 = 0
	// RespMsg 给客户端 的 响应
	RespMsg uint8 = 1
	// PushMsg 给客户端 的 推送, 需要ack确认
	PushMsg uint8 = 2
)

// sendInfo
type sendInfo struct {
	// 发送 的 类型
	t uint8
	// 发送的值
	data []byte
}

// CliConn 客户连接
type CliConn struct {
	s     *Server
	conn  *websocket.Conn
	ch    chan *sendInfo
	ack   chan bool
	state int32 // conn 状态，0:nil, 1:running, 2:closing, 3:closed

	name  string // 客户名称
	kind  uint8  // 客户端类型
	flags uint32 // 客户端 属性
}

func (c *CliConn) isRunning() bool {
	return atomic.LoadInt32(&c.state) == RunningConnState
}

func (c *CliConn) pushString(msg string) bool {
	// 再一次 判断 连接 是否 断开
	if !c.isRunning() {
		return false
	}
	info := &sendInfo{
		t:    PushMsg,
		data: []byte(msg),
	}
	c.ch <- info
	return true
}

func (c *CliConn) pushBytes(msg []byte, variable bool) bool {
	// 再一次 判断 连接 是否 断开
	if !c.isRunning() {
		return false
	}
	info := &sendInfo{
		t:    PushMsg,
		data: msg,
	}

	if variable { // db 查询时，共用一个[]byte地址，会变，所以这里 拷贝一份
		// 而在 insert 方法里  是直接new 的 所以不会变
		info.data = make([]byte, len(msg))
		copy(info.data, msg)
	}
	c.ch <- info
	return true
}

func (c *CliConn) logout(tip string) {
	if !c.isRunning() {
		return
	}
	// 设置 正在关闭
	atomic.StoreInt32(&c.state, ClosingConnState)

	info := &sendInfo{
		t: LogoutMsg,
	}
	info.data = websocket.FormatCloseMessage(websocket.CloseNormalClosure, tip)
	c.ch <- info
}

// 发送 响应
func (c *CliConn) sendResp(body []byte) {
	if !c.isRunning() {
		return
	}

	info := &sendInfo{
		t:    RespMsg,
		data: body,
	}

	c.ch <- info
}

//
func cliServe(s *Server, conn *websocket.Conn) {
	c := &CliConn{
		s:     s,
		conn:  conn,
		ch:    make(chan *sendInfo, 20),
		ack:   make(chan bool),
		state: 0,
	}
	// 新建一个 gorou 用来 并发 发送 消息
	go c.handle()
	c.Serve()
	close(c.ch)
}

// sendMsgAndWaitAck 发送 消息到客户端 并且 等待 ack 回复，
func (c *CliConn) sendMsgAndWaitAck(data []byte) {
	// 三次 没接收到ack 就 当 离线了
	// 推送 到 客户端 是 串行的，push-ack，
	for i := 0; i < 3; i++ {
		err := c.conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			c.s.log.Error("WriteMessage err:", err)
			//
			c.logout("errWrite")
			return
		}

		select {
		case <-c.ack:
			// 收到 ack 回复 后 就直接 返回
			return
		case <-time.After(time.Second):
			// 超时后 就 重发
			continue
		}
	}
	// 运行到这里 就表示 发送失败了
	c.logout("notAck")
}

// handle 长连接 里的 待发送消息 队列
func (c *CliConn) handle() {
	for {
		msg := <-c.ch
		if msg == nil {
			// 发送空 来 判断 是否 退出
			break
		}
		// logout 退出
		if LogoutMsg == msg.t {
			// 主动 关闭 然后 read 里 就会 报错
			_ = c.conn.WriteControl(websocket.CloseMessage, msg.data, time.Time{})
			// 这里 不推出 还是由 read 里 报错 在 <- nil 退出
			continue
		} else if RespMsg == msg.t {
			// resp 回复 客户端的请求 resp
			err := c.conn.WriteMessage(websocket.TextMessage, msg.data)
			if err != nil {
				c.s.log.Error("handleSend.Write err:", err)
				//
				c.logout("errResp")
			}
		} else {
			// 其余 的 就是需要 ack 的了
			c.sendMsgAndWaitAck(msg.data)
		}
	}
	c.s.log.Debug("the Send close:", c.conn.RemoteAddr().String())
}

// Serve 管理 客户端接入 的 长连接
func (c *CliConn) Serve() {
	c.s.log.Debug("the websocket connect:", c.conn.RemoteAddr().String())
	for {
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			// 断开后 第一时间 设置 关闭状态
			atomic.StoreInt32(&c.state, 3)

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
			c.logout("notJson")
			continue
		}

		fun := j.GetString("fun")
		if len(fun) == 0 {
			c.s.log.Error("异常连接, err json", j)
			c.logout("errJson")
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
	flags := j.GetUint32("flags")

	// 添加 验证 用户

	if len(c.name) > 0 { // 已经登录过了
		return
	}
	c.name = name
	c.kind = 0xFF & uint8(flags>>24) // 高8位
	c.flags = 0xFFFFFF & flags       // 低24位

	c.s.setAndKickUser(c)
	// 设置 运行状态
	atomic.StoreInt32(&c.state, RunningConnState)

	// 添加 返回信息
	c.sendResp([]byte(`{"ret": "login", "text": "success"}`))

	// 异步 查询 消息
	c.s.job.Add(func() {
		// 读取 db 里 未读 的 消息
		//delList := list.New()
		preKey := makeDbMsgPre(c.name)
		ite := c.s.db.NewIterator(util.BytesPrefix(preKey), nil)
		for ite.Next() {
			c.pushBytes(ite.Value(), true)
			err := c.s.db.Delete(ite.Key(), nil)
			if err != nil {
				c.s.log.Error("db.Write err:", err)
			}
			//buf := make([]byte, len(ite.Key()))
			//copy(buf, ite.Key())
			//delList.PushBack(buf)
		}
		ite.Release()
		/*
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

		*/
	})

}
