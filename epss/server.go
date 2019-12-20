package epss

import (
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

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) setAndKickUser(name string, kind uint8, p *CliConn) {
	key := makeCliKey(name, kind)
	// 先 判断 该用户 是否已经登录了
	if old, ex := s.mapConns.Load(key); ex {
		// 登录了 就先踢他下线
		if conn, ok := old.(*CliConn); ok && !conn.closed {
			conn.logout("kickout")
		}
	}
	// 设置 在线 了
	s.mapConns.Store(key, p)
}

func (s *Server) getPusher(name string, kind uint8) (p Pusher) {
	key := makeCliKey(name, kind)

	if v, ex := s.mapConns.Load(key); ex {
		if conn, ok := v.(*CliConn); ok && !conn.closed {
			p = conn
		}
	}

	return
}
