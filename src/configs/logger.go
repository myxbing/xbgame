package configs

import (
	"io"
	"log"
	"os"
)

type DoubleLogger struct {
	out1, out2 io.Writer
}

// 实现io.Writer接口
func (r *DoubleLogger) Write(p []byte) (int, error) {
	var n int
	var e error

	if r.out1 != nil {
		n, e = r.out1.Write(p)
	}

	if r.out2 != nil {
		n, e = r.out2.Write(p)
	}

	return n, e
}

// 设置我们的日志记录器
func StartLogger(logfile string) {
	bl := []byte(logfile)

	var err error
	var f *os.File

	if bl[0] == '/' { // start with slash, just open
		f, err = os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	} else {
		path := os.Getenv("GOPATH") + "/" + logfile
		f, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	}

	if err != nil {
		log.Println("cannot open logfile %v\n", err)
		os.Exit(-1)
	}

	var r DoubleLogger

	switch String("log.mode") {
	case "both":
		r.out1 = os.Stdout
		r.out2 = f
	case "file":
		r.out2 = f
	}
	log.SetOutput(&r)
}
