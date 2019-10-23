package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ijesonchen/glog"

	"github.com/ijesonchen/filededup/utils"
)

func setglog(toConsole bool, logdir, level string) {
	var err error
	glog.AlsoToStderr(toConsole)
	if len(logdir) > 0 {
		if err = os.MkdirAll(logdir, 0777); err != nil {
			glog.Error("os.MkdirAll error: ", err)
			return
		}
		glog.SetLogDir(logdir)
	}
	if len(level) > 0 {
		if err = glog.SetLevel(level); err != nil {
			glog.Fatal("SetLevel error: ", err)
			return
		}
	}
}

func main() {
	var err error
	var dirTotal, fileTotal int32
	defer glog.Flush()
	defer glog.Infof("main exit.")
	defer func() { glog.Infof("total %d dirs, %d files", dirTotal, fileTotal) }()

	if len(os.Args) <= 1 {
		glog.Error("no dir speified.")
		return
	}

	setglog(true, "log_dir", "INFO")

	initDirs := os.Args[1:]

	flag.Parse()
	for i, a := range os.Args {
		fmt.Println(i, a)
	}

	cfg := ReadConfig("filedup.json")
	setglog(cfg.Log2Console, cfg.LogDir, cfg.LogLevel)

	// log flusher thread
	ticker := time.NewTicker(cfg.LogFlushSec * time.Second)
	go func() {
		for range ticker.C {
			glog.Flush()
		}
	}()

	chDir := make(chan string)
	chFile := make(chan string)

	var dirJobLeft, fileJobLeft int32
	dirJobLeft = int32(len(initDirs))
	dirTotal = 0 // do not count root dir

	go func() {
		for _, d := range initDirs {
			chDir <- d
		}
	}()

	var wgWalker, wgHasher, wgJob sync.WaitGroup

	var hasherStarter sync.Once
	startHasher := func() {
		wgHasher.Add(cfg.HasherThreads)
		for i := 0; i < cfg.HasherThreads; i++ {
			go func(i int) {
				for fn := range chFile {
					h, e := fnv64File(fn, cfg.Byte2Hash)
					currentFile := atomic.AddInt32(&fileJobLeft, -1)
					currentDir := atomic.LoadInt32(&dirJobLeft)
					if e != nil {
						glog.Errorf("hashFile %q error %v", fn, e)
					} else {
						glog.Infof("--> %q : %d", fn, h)
					}
					// if finished...
					if currentFile == 0 && currentDir == 0 {
						close(chFile)
						break
					}
				}
				wgHasher.Done()
				glog.Infof("hasher thread %d stopped.", i)
			}(i)
		}
	}

	wgWalker.Add(cfg.WalkerThreads)
	for i := 0; i < cfg.WalkerThreads; i++ {
		go func(i int) {
			// get next dir
			for dir := range chDir {
				// process dir
				dirs, files, e := walkDir(dir)
				nDir := int32(len(dirs))
				nFile := int32(len(files))
				// cnt job first
				atomic.AddInt32(&dirTotal, nDir)
				current := atomic.AddInt32(&dirJobLeft, nDir-1)
				if err != nil {
					glog.Errorf("walk dir %q error: %v", dir, e)
					// if finished...
					if current == 0 {
						close(chDir)
						break
					}
					continue
				}
				atomic.AddInt32(&fileTotal, nFile)
				atomic.AddInt32(&fileJobLeft, nFile)

				// send jobs
				go func(dirs []string) {
					for _, i := range dirs {
						chDir <- i
					}
				}(dirs)
				go func(files []string) {
					for _, i := range files {
						// make sure at least on file
						hasherStarter.Do(startHasher)
						chFile <- i
					}
				}(files)

				// if finished...
				if current == 0 {
					close(chDir)
					break
				}
			}
			wgWalker.Done()
			glog.Warningf("walker thread %d stopped.", i)
		}(i)
	}

	// job waiter
	wgJob.Add(1)
	go func() {
		wgWalker.Wait()
		glog.Warning("all walker thread stopped.")

		wgHasher.Wait()
		glog.Warning("all hasher thread stopped.")
		wgJob.Done()
	}()

	// wait all worker threads
	glog.Info("wait all job threads.")
	wgJob.Wait()

	// for concurrent control
	cc := utils.NewConcurentControl(20, 100, 10)
	cc.Enter()
	cc.Leave()
}
