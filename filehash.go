package main

import (
	"crypto/sha1"
	"fmt"
	"hash/fnv"
	"os"

	"github.com/ijesonchen/glog"
)

func readData(fn string, byte2hash int) (data []byte, err error) {
	if len(fn) == 0 || byte2hash <= 0 {
		err = fmt.Errorf("invalid parameter: %q, %d", fn, byte2hash)
		return
	}
	st, err := os.Stat(fn)
	if err != nil {
		glog.Errorf("Stat file %s error %v", fn, err)
		return
	}
	if st.IsDir() {
		err = fmt.Errorf("hash dir")
		return
	}
	if st.Size() < int64(byte2hash) {
		err = fmt.Errorf("to small: %d to %d", st.Size(), byte2hash)
		return
	}

	f, err := os.Open(fn)
	if err != nil {
		glog.Errorf("open file %s error %v", fn, err)
		return
	}

	data = make([]byte, byte2hash)
	n, err := f.Read(data)
	if n != byte2hash {
		err = fmt.Errorf("read byte error: expect %d got %d", byte2hash, n)
		return
	}

	return
}

func fnv64File(fn string, byte2hash int) (fp uint64, err error) {
	glog.Infof("ENTER fnv64File %q, %d", fn, byte2hash)
	defer func() {
		if err != nil {
			glog.Errorf("fnv64File %q error: %v", fn, err)
		}
		glog.Infof("LEAVE fnv64File %q", fn)
	}()

	data, err := readData(fn, byte2hash)
	h := fnv.New64a()
	h.Write(data)
	fp = h.Sum64()
	if fp == 0 {
		glog.Warningf("fnv64File %q got 0 hash", fn)
	}
	return
}

func sha1File(fn string, byte2hash int) (hash []byte, err error) {
	glog.Infof("ENTER sha1File %q, %d", fn, byte2hash)
	defer func() {
		if err != nil {
			glog.Errorf("sha1File %q error: %v", fn, err)
		}
		glog.Infof("LEAVE sha1File %q", fn)
	}()

	data, err := readData(fn, byte2hash)

	h := sha1.New()
	h.Write(data)
	hash = h.Sum(nil)
	return
}

// file size should equal. only compare first byte2Comp bit
func compFileHeadContent(fn1, fn2 string, byte2Comp int) (same bool, err error) {
	glog.Infof("ENTER compFileHeadContent %q %q, %d", fn1, fn2, byte2Comp)
	defer func() {
		if err != nil {
			glog.Errorf("compFileHeadContent %q %q error: %v", fn1, fn2, err)
		}
		glog.Infof("LEAVE compFileHeadContent %q %q", fn1, fn2)
	}()

	d1, err := readData(fn1, byte2Comp)
	if err != nil {
		glog.Errorf("read %q error %v", fn1, err)
	}
	d2, err := readData(fn2, byte2Comp)
	if err != nil {
		glog.Errorf("read %q error %v", fn2, err)
	}

	for i := 0; i < byte2Comp; i++ {
		if d1[i] != d2[i] {
			return
		}
	}

	same = true
	return
}
