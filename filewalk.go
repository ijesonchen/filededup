package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ijesonchen/glog"
)

func isDir(fn string) bool {
	if len(fn) == 0 {
		glog.Errorf("isDir invalid parameter: %q", fn)
		return false
	}
	st, err := os.Stat(fn)
	if err != nil {
		glog.Errorf("isDir Stat file %s error %v", fn, err)
		return false
	}
	return st.IsDir()
}

func walkDir(fn string) (dirs, files []string, err error) {
	glog.Infof("ENTER walkDir %q", fn)
	defer func() {
		if err != nil {
			glog.Errorf("walkDir %q error: %v", fn, err)
		}
		glog.Infof("LEAVE walkDir %q", fn)
	}()

	if !isDir(fn) {
		err = fmt.Errorf("not dir")
	}

	infos, err := ioutil.ReadDir(fn)
	if err != nil {
		return
	}

	for _, info := range infos {
		name := filepath.Join(fn, info.Name())
		if info.IsDir() {
			dirs = append(dirs, name)
		}
		if info.Mode().IsRegular() {
			files = append(files, name)
		}
	}

	return
}
