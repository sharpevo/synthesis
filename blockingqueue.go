package blockingqueue

import (
	"conditionvariable"
	"fmt"
	"sync"
)

type BlockingQueue struct {
	sync.Mutex
	itemList      []interface{}
	pushCondition conditionvariable.ChannelCondition
	termCondition conditionvariable.ChannelCondition
}

func NewBlockingQueue() *BlockingQueue {
	return &BlockingQueue{
		itemList:      []interface{}{},
		pushCondition: conditionvariable.NewChannelCondition(),
		termCondition: conditionvariable.NewChannelCondition(),
	}
}

func (b *BlockingQueue) Reset() {
	b.Lock()
	defer b.Unlock()
	b.itemList = []interface{}{}
	b.termCondition.Broadcast()
}

func (b *BlockingQueue) Push(item interface{}) {
	b.Lock()
	defer b.Unlock()
	b.itemList = append(b.itemList, item)
	//fmt.Println("pushing", item, b.itemList)
	b.pushCondition.Broadcast()
}

func (b *BlockingQueue) Pop() (interface{}, error) {
	b.Lock()
	defer b.Unlock()
	for len(b.itemList) == 0 {
		b.Unlock()
		select {
		case <-b.pushCondition.Wait():
			b.Lock()
			continue
		case <-b.termCondition.Wait():
			b.Lock()
			return nil, fmt.Errorf("queue terminated")
		}
	}
	item := b.itemList[0]
	b.itemList = b.itemList[1:]
	//fmt.Println("popping", item, b.itemList)
	return item, nil
}

func (b *BlockingQueue) Length() int {
	b.Lock()
	defer b.Unlock()
	return len(b.itemList)
}

func (b *BlockingQueue) Get(index int) (interface{}, error) {
	b.Lock()
	defer b.Unlock()
	if index < 0 || index > len(b.itemList) {
		return nil, fmt.Errorf("invalid index")
	}
	return b.itemList[index], nil
}

func (b *BlockingQueue) Append(item interface{}) {
	b.Lock()
	defer b.Unlock()
	b.itemList = append(b.itemList, item)
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
