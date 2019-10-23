package utils

import (
	"time"

	"github.com/ijesonchen/glog"
)

// ConcurentControl ...
type ConcurentControl struct {
	maxConCurrent int
	maxWait       time.Duration
	sleepTime     time.Duration
	tokenChan     chan struct{}
}

// NewConcurentControl ...
func NewConcurentControl(maxConCurrent, maxWaitMs, sleepMs int) (cc *ConcurentControl) {
	cc = &ConcurentControl{
		maxConCurrent: maxConCurrent,
		maxWait:       time.Millisecond * time.Duration(maxWaitMs),
		sleepTime:     time.Millisecond * time.Duration(sleepMs),
		tokenChan:     make(chan struct{}, maxConCurrent),
	}
	return
}

// Enter ...
func (cc *ConcurentControl) Enter() (ok bool) {
	t0 := time.Now()
	for {
		select {
		case cc.tokenChan <- struct{}{}:
			return true
		default:
			time.Sleep(cc.sleepTime)
			if time.Now().Sub(t0) > cc.maxWait {
				return false
			}
		}
	}
}

// Leave ...
func (cc *ConcurentControl) Leave() {
	select {
	case <-cc.tokenChan:
	default:
		glog.Errorf("ConcurentControl token error")
	}
}
