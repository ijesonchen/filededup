package db

import (
	"fmt"

	"github.com/ijesonchen/glog"

	"github.com/globalsign/mgo"
)

// NewClient .
func NewClient(host string, port uint16, dbName, collName string) (client IClient, err error) {

	// [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
	url := fmt.Sprintf("mongodb://%s:%d/%s", host, port, dbName)
	sess, err := mgo.Dial(url)
	if err != nil {
		glog.Errorf("connect to %s error: %v", url, err)
		return
	}
	coll := sess.DB("").C(collName)

	client = &Client{
		url:      url,
		collName: collName,
		sess:     sess,
		coll:     coll,
	}

	return
}

// Client .
type Client struct {
	url      string
	collName string
	sess     *mgo.Session
	coll     *mgo.Collection
}

// Insert .
func (c *Client) Insert(d *FileInfo) (err error) {
	return
}

// FindName .
func (c *Client) FindName(name string) (files []FileInfo, err error) {
	return
}

// LikeName .
func (c *Client) LikeName(name string) (files []FileInfo, err error) {
	return
}

// Delete .
func (c *Client) Delete(id string) (err error) {
	return
}

// List .
func (c *Client) List(count, skip int64, order string) (files []FileInfo, err error) {
	return
}
