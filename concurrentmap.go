package concurrentmap

import (
	"fmt"
	"sync"
)

type ConcurrentMap struct {
	sync.RWMutex
	m map[string]interface{}
}

type Item struct {
	Key   string
	Value interface{}
}

func NewConcurrentMap(cmaps ...*ConcurrentMap) *ConcurrentMap {
	newMap := make(map[string]interface{})
	if len(cmaps) != 0 {
		for _, cmap := range cmaps {
			for item := range cmap.Iter() {
				newMap[item.Key] = item.Value
			}
		}
	}
	return &ConcurrentMap{
		m: newMap,
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

func (c *ConcurrentMap) Del(key string) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.m[key]; !ok {
		return fmt.Errorf("failed to delete map entry due to invalid key %s", key)
	}
	delete(c.m, key)
	return nil
}

func (c *ConcurrentMap) Replace(ori interface{}, value interface{}) (key string, err error) {
	c.Lock()
	defer c.Unlock()
	for k, v := range c.m {
		if ori == v {
			c.m[k] = value
			return k, nil
		}
	}
	return key, fmt.Errorf("%v not existed", ori)
}

func (c *ConcurrentMap) String() string {
	c.RLock()
	defer c.RUnlock()
	return fmt.Sprintf("%#v", c.m)
}

func (c *ConcurrentMap) Iter() <-chan Item {
	itemc := make(chan Item)
	go func() {
		defer close(itemc)
		c.Lock() // more secure than RLock
		defer c.Unlock()
		for k, v := range c.m {
			itemc <- Item{k, v}
		}
	}()
	return itemc
}
