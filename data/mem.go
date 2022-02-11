package data

import (
	"fmt"
	"math/rand"
	"sync"
)

type MemoryStateManager struct {
	data map[string]*GameState
	subs map[string]map[chan string]bool

	dataLock, subLock sync.RWMutex
}

func NewMemoryStateManager() *MemoryStateManager {
	return &MemoryStateManager{
		data: map[string]*GameState{},
		subs: map[string]map[chan string]bool{},
	}
}

func (r *MemoryStateManager) genKey() string {
	var num int
	var key string
	for {
		num = rand.Intn(1048576)
		key = fmt.Sprintf("%05X", num)

		r.subLock.RLock()
		if _, ok := r.subs[key]; !ok {
			return key
		}
		r.subLock.RUnlock()
	}
}
func (r *MemoryStateManager) Get(key string) *GameState {
	r.dataLock.RLock()
	defer r.dataLock.RUnlock()

	return r.data[key]
}
func (r *MemoryStateManager) Set(key string, state *GameState) bool {
	r.dataLock.Lock()
	defer r.dataLock.Unlock()

	r.data[key] = state
	return true
}
func (r *MemoryStateManager) Del(key string) bool {
	r.dataLock.Lock()
	defer r.dataLock.Unlock()

	delete(r.data, key)
	return true
}

func (r *MemoryStateManager) Sub(channel string) chan string {
	r.subLock.Lock()
	defer r.subLock.Unlock()

	ch := make(chan string)
	if r.subs[channel] == nil {
		r.subs[channel] = map[chan string]bool{}
		r.Set(channel, new(GameState))
	}
	r.subs[channel][ch] = true

	return ch
}
func (r *MemoryStateManager) Unsub(channel string, in chan string) bool {
	r.subLock.Lock()
	defer r.subLock.Unlock()

	close(in)
	delete(r.subs[channel], in)
	if len(r.subs[channel]) == 0 {
		r.Del(channel)
		delete(r.subs, channel)
	}

	return true
}
func (r *MemoryStateManager) Pub(msg, channel string, in chan string) bool {
	r.subLock.RLock()
	defer r.subLock.RUnlock()

	for ch := range r.subs[channel] {
		if ch != in {
			ch <- msg
		}
	}
	return true
}
