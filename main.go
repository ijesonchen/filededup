package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
)

func main() {
	var err error
	defer glog.Exit("main exit.")

	if len(os.Args) <= 1 {
		glog.Error("no dir speified.")
		return
	}

	initDirs := os.Args[1:]

	flag.Parse()
	for i, a := range os.Args {
		fmt.Println(i, a)
	}

	glog.AlsoToStderr(true)
	logdir := "log"
	if err = os.MkdirAll(logdir, 0777); err != nil {
		glog.Error("os.MkdirAll error: ", err)
		return
	}
	glog.SetLogDir(logdir)
	loglevel := "INFO"
	if err = glog.SetLevel(loglevel); err != nil {
		glog.Fatal("SetLevel error: ", err)
		return
	}

	cfg := ReadConfig("filedup.json")

	chDir := make(chan string)
	chFile := make(chan string)

	go func() {
		for _, d := range initDirs {
			chDir <- d
		}
	}()

	for i := 0; i < cfg.WalkerThreads; i++ {
		go func() {
			for dir := range chDir {
				dirs, files, e := walkDir(dir)
				if err != nil {
					glog.Errorf("walk dir %q error: %v", dir, e)
					continue
				}
				go func() {
					for _, i := range dirs {
						chDir <- i
					}
				}()
				go func() {
					for _, i := range files {
						chFile <- i
					}
				}()
			}
		}()
	}

	for i := 0; i < cfg.HasherThreads; i++ {
		go func() {
			for fn := range chFile {
				h, e := hashFile(fn, cfg.Byte2Hash)
				if e != nil {
					glog.Errorf("hashFile %q error %v", fn, e)
				}
				glog.Warningf("--> %q : %d", fn, h)
			}
		}()
	}

	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			glog.Flush()
		}
	}()

	barrer := make(chan bool)
	<-barrer
}
