package epss

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"net/http"
	"sync"
	//_ "github.com/mattn/go-sqlite3"
)

// 日志 接口
type logger interface {
	Debug(v ...interface{})
	Error(v ...interface{})
}

// Server epss
type Server struct {
	log logger

	// 管理 websocket 连接
	mapConns sync.Map
	job      jobPool

	db *leveldb.DB
}

// New a app
func New(log logger) *Server {
	s := &Server{}

	if log != nil {
		s.log = log
	} else {
		s.log = d
	}

	s.job.limit = 100
	return s
}

// ListenAndServe run
func (s *Server) ListenAndServe(addr string) error {
	var err error

	s.db, err = leveldb.OpenFile("db", nil)
	if err != nil {
		s.log.Error("leveldb open err:", err)
		return err
	}
	var fla int = int(KindDocsy)
	fla = (fla << 24) + int(KindDocsy)
	fmt.Println(fla)
	makeDbMsgKey("namm", nil, 2, "12313")
	//bb, err := s.db.Get([]byte("wuxb:1"), nil)
	//fmt.Println(bb)

	// db, err := sql.Open("sqlite3", ".\\msg.db")
	// if err != nil {

	// }
	// defer db.Close()

	// sqlStmt := `
	// create table foo (id integer not null primary key, name text);
	// delete from foo;
	// `
	// _, err = db.Exec(sqlStmt)
	// if err != nil {
	// 	log.Printf("%q: %s\n", err, sqlStmt)
	// 	return nil
	// }

	http.HandleFunc("/state/online", s.handleOnline)
	http.HandleFunc("/msg/inserts", s.handleInserts)
	http.HandleFunc("/msg/insert", s.handleInsert)
	http.HandleFunc("/msg/remove", s.handleRemove)
	http.HandleFunc("/conn/cli", s.handleClient)

	return http.ListenAndServe(addr, nil)
}

func (s *Server) setAndKickUser(p *CliConn) {
	// 先 判断 该用户 是否已经登录了
	v, ex := s.mapConns.Load(p.name)
	if !ex {
		// 设置 在线 了
		conn := []*CliConn{p}
		s.mapConns.Store(p.name, conn)
		return
	}

	// 查询到
	if conns, ok := v.([]*CliConn); ok {
		for i, v := range conns {
			// 登录了 就先踢他下线
			if v.kind == p.kind {
				if v.isRunning() {
					v.logout("kickout")
				}
				conns[i] = p
				return
			}
		}
		// 没有 在内 就 append
		conns = append(conns, p)
	}
}

// 方法封装
func (s *Server) pushMsgTo(name string, src uint8, msgBytes []byte, id string) {
	// name key 查询 所有 在线 连接 推送过去 只要 一个成功了 后面就 不 存储db
	suc := 0
	if v, ex := s.mapConns.Load(name); ex {
		if conns, ok := v.([]*CliConn); ok {
			for _, v := range conns {
				if v.isRunning() {
					if v.pushBytes(msgBytes, false) {
						suc++
					}
				}
			}
		}
	}

	if suc > 0 {
		return
	}

	// 不在线 或 push失败 		 数据 保存 的 数据库
	dbKey := makeDbMsgKey(name, nil, src, id)
	if err := s.db.Put(dbKey, msgBytes, nil); err != nil {
		s.log.Error("dbPut err:", err)
	}
}
