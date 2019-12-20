package main

import (
	"epssd/epss"
	"fmt"
	"path/filepath"
	"time"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
)

var theLog = glog.DefaultLogger()

func main() {
	var err error
	// 自身的位置
	selfDir := gfile.SelfDir()
	// 日志 对象
	logPath := filepath.Join(selfDir, "logs")
	fmt.Println("the log save at ", logPath)
	_ = theLog.SetPath(logPath)
	theLog.SetAsync(true)

	// new
	ser := epss.New(theLog)
	err = ser.ListenAndServe("localhost:8345")
	if err != nil {
		theLog.Error("Listen err:", err)
		// 保证 async 日志 写完
		time.Sleep(time.Second)
	}
}
