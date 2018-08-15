package blockingqueue

import (
	"sync"
)

type BlockingQueue struct {
	sync.RWMutex
	itemList []interface{}
	popc     chan interface{}
}

func NewBlockingQueue() *BlockingQueue {
	return &BlockingQueue{
		itemList: []interface{}{},
		popc:     make(chan interface{}),
	}
}

func (b *BlockingQueue) Reset() {
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

func (b *BlockingQueue) Pop() interface{} {
	if len(b.itemList) <= 0 {
		<-b.popc
	}
	b.Lock()
	defer b.Unlock()
	item := b.itemList[0]
	b.itemList = b.itemList[1:]
	return item
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
