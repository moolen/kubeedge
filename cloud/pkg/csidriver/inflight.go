package csidriver

import (
	"sync"
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
	defer i.mux.Unlock()
	if _, ok := i.lookup[key]; ok {
		return false
	}
	i.lookup[key] = true
	return true
}
func (i *inFlight) Delete(key string) {
	i.mux.Lock()
	defer i.mux.Unlock()
	delete(i.lookup, key)
}
