package epss

import (
	"epssd/json"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
)

func (s *Server) getPostBodyJSON(w http.ResponseWriter, req *http.Request) *json.JSON {
	if "POST" != req.Method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}

	var err error
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		s.log.Error("ReadBody err:", err)
		return nil
	}

	js, err := json.LoadBytes(body)
	if err != nil {
		s.log.Error("LoadJson err:", err)
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}
	return js
}

// 插入 提供 一次 单人 单客户端 的 消息
func (s *Server) handleInsert(w http.ResponseWriter, req *http.Request) {
	js := s.getPostBodyJSON(w, req)
	if js == nil {
		return
	}

	to := js.GetString("to")       	// 用户名
	src := js.GetUint8("src")    	// 消息 来自 哪个 应用系统
	msg := js.GetJSON("msg") 		// 消息体 json 必须 要有 id
	if msg == nil {
		s.log.Error("not msg")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	msg.Add("ret","push")
	// id 作为消息的唯一标志，用来 服务端和客户端 来 定位具体的某个消息
	id := msg.GetString("id","0000")
	msgBytes := msg.Marshal()

	// 添加 异步任务
	fn := func() {
		s.pushMsgTo(to, src, msgBytes, id)
	}
	s.job.Add(fn)
	return
}

// 提供 一次 插入 多人 多客户端 的 消息
func (s *Server) handleInserts(w http.ResponseWriter, req *http.Request) {
	// 解析出 消息 模板
	js := s.getPostBodyJSON(w, req)
	if js == nil {
		return
	}
	var tos []string
	toJ := js.GetJSON("tos") 	// 用户名 多个
	if toJ != nil {
		tos = toJ.ToStrings()
	}
	src := js.GetUint8("src")    // 消息 来自 哪个 应用系统
	msg := js.GetJSON("msg") // 消息体 json 必须 要有 id
	if msg == nil {
		s.log.Error("not msg")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := msg.GetString("id", "0000")
	msgBytes := msg.Marshal()

	// 添加 异步任务
	fn := func() {
		for _, to := range tos {
			s.pushMsgTo(to, src, msgBytes, id)
		}
	}
	s.job.Add(fn)
	return
}



// 消息 没有及时 推送出去 会 存储在 数据库，这时 是 可以 移除的
// 就算推送到了终端上，还可以继续 推送移除指令
func (s *Server) handleRemove(w http.ResponseWriter, req *http.Request) {
	js := s.getPostBodyJSON(w, req)
	if js == nil {
		return
	}
}

// 获取 在线用户
func (s *Server) handleOnline(w http.ResponseWriter, req *http.Request) {
	_, _ = w.Write([]byte("Coding ..."))
}

//websocket 升级
var _upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (s *Server) handleClient(w http.ResponseWriter, req *http.Request) {
	conn, err := _upgrader.Upgrade(w, req, nil)
	if err != nil {
		s.log.Error("Upgrade err:", err)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	cliServe(s, conn)
}
