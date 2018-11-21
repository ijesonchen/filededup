package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ijesonchen/glog"
)

// Config .
type Config struct {
	WalkerThreads int // threads to search folders
	HasherThreads int // threads to hash files
	MinValidByte  int // minimux file size to process
	Byte2Hash     int // less than MinValidByte. only hash from part of file with hash32.
}

// ReadConfig .
func ReadConfig(fn string) *Config {
	cfg := &Config{}
	err := cfg.read(fn)
	if err != nil {
		glog.Fatalf("read config %q error: %v", fn, err)
	}

	glog.Infof("config %+v", *cfg)
	return cfg
}

// IsValid check parameter
func (c *Config) isValid() (valid bool) {
	if c.WalkerThreads <= 0 || c.WalkerThreads > 100 {
		c.WalkerThreads = 10
		glog.Info("WalkerThreads set to 10")
	}
	if c.HasherThreads <= 0 || c.HasherThreads > 100 {
		c.HasherThreads = 10
		glog.Info("WalkerThreads set to 10")
	}
	if c.MinValidByte == 0 {
		c.HasherThreads = 1024 * 1024 * 1024 * 50
		glog.Info("WalkerThreads set to 50M")
	}
	if c.Byte2Hash == 0 {
		c.HasherThreads = 1024 * 1024 * 1024 * 5
		glog.Info("WalkerThreads set to 5M")
	}
	if c.Byte2Hash > c.MinValidByte {
		glog.Errorf("Byte2Hash %d > MinValidByte %d", c.Byte2Hash, c.MinValidByte)
		return false
	}
	return true
}

func (c *Config) read(fn string) (err error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		glog.Fatalf("open file %s error: %v", fn, err)
		return
	}
	if err = json.Unmarshal(data, &c); err != nil {
		glog.Fatalf("parse file %s error: %v", fn, err)
		return
	}
	if !c.isValid() {
		glog.Fatalf("config file %s not valid", fn)
		return
	}
	return
}
