package csidriver

import (
	"sync"

	"k8s.io/klog/v2"
)

type inFlight struct {
	mux    *sync.Mutex
	lookup map[string]bool
}

func newInFlight() *inFlight {
	return &inFlight{
		mux:    &sync.Mutex{},
		lookup: map[string]bool{},
	}
}

func (i *inFlight) Insert(key string) bool {
	i.mux.Lock()
	defer i.mux.Lock()
	if _, ok := i.lookup[key]; ok {
		klog.Infof("inflight: already locked %s", key)
		return false
	}
	klog.Infof("inflight: locking %s", key)
	i.lookup[key] = true
	return true
}
func (i *inFlight) Delete(key string) {
	i.mux.Lock()
	defer i.mux.Unlock()
	klog.Infof("inflight: unlocking %s", key)
	delete(i.lookup, key)
}
