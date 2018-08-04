package concurrentmap

import (
	"sync"
)

type ConcurrentMap struct {
	sync.RWMutex
	m map[string]interface{}
}

func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		m: make(map[string]interface{}),
	}
}

func (c *ConcurrentMap) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()
	value, ok := c.m[key]
	return value, ok
}

func (c *ConcurrentMap) Set(key string, value interface{}) interface{} {
	c.Lock()
	defer c.Unlock()
	c.m[key] = value
	return value
}
