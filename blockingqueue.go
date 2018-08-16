package blockingqueue

import (
	"fmt"
	"sync"
)

type BlockingQueue struct {
	sync.RWMutex
	itemList   []interface{}
	popc       chan interface{}
	terminatec chan interface{}
}

func NewBlockingQueue() *BlockingQueue {
	return &BlockingQueue{
		itemList:   []interface{}{},
		popc:       make(chan interface{}),
		terminatec: make(chan interface{}),
	}
}

func (b *BlockingQueue) Reset() {
	select {
	case b.terminatec <- true:
	default:
	}
	b.itemList = []interface{}{}
}

func (b *BlockingQueue) Push(item interface{}) {
	b.Lock()
	defer b.Unlock()
	b.itemList = append(b.itemList, item)
	if len(b.itemList) == 1 {
		select {
		case b.popc <- true:
		default:
		}
	}
}

func (b *BlockingQueue) Pop() (interface{}, error) {
	if len(b.itemList) <= 0 {
		for {
			select {
			case <-b.popc:
				return b.pop()
			case <-b.terminatec:
				return nil, fmt.Errorf("queue terminated")
			}
		}
	}
	return b.pop()
}

func (b *BlockingQueue) pop() (interface{}, error) {
	b.Lock()
	defer b.Unlock()
	item := b.itemList[0]
	b.itemList = b.itemList[1:]
	return item, nil
}

type Item struct {
	Index int
	Value interface{}
}

func (b *BlockingQueue) Iter() <-chan Item {
	itemc := make(chan Item)
	go func() {
		defer close(itemc)
		b.Lock()
		defer b.Unlock()
		for k, v := range b.itemList {
			itemc <- Item{k, v}
		}
	}()
	return itemc
}
