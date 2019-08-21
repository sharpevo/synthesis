package blockingqueue

import (
	"conditionvariable"
	"fmt"
	"log"
	"runtime/debug"
	"sync"
)

type BlockingQueue struct {
	lock          sync.Mutex
	itemList      []interface{}
	Cursor        chan int
	pushCondition conditionvariable.ChannelCondition
	termCondition conditionvariable.ChannelCondition
}

func NewBlockingQueue() *BlockingQueue {
	b := &BlockingQueue{}
	b.itemList = []interface{}{}
	b.Cursor = make(chan int)
	b.pushCondition = conditionvariable.NewChannelCondition()
	b.termCondition = conditionvariable.NewChannelCondition()
	return b
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
	if index < 0 || index >= len(b.itemList) {
		return nil, fmt.Errorf("invalid index")
	}
	item := b.itemList[index]
	return item, nil
}

func (b *BlockingQueue) GetLockless(index int) (interface{}, error) {
	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf(
				"Panic: %s\n%s",
				err, debug.Stack(),
			)
			log.Println("---------- panic ----------")
			log.Fatal(msg)
		}
	}()
	if index < 0 || index >= len(b.itemList) {
		return nil, fmt.Errorf("invalid index")
	}
	item := b.itemList[index]
	return item, nil
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

func (b *BlockingQueue) Next() <-chan Item {
	itemc := make(chan Item)
	go func() {
		defer close(itemc)
		b.Lock()
		defer b.Unlock()
		var cursor int
		for cursor < len(b.itemList) {
			cursor = <-b.Cursor
			fmt.Println("CURSOR", cursor, len(b.itemList))
			fmt.Printf("CURSOR queue: %#v\n", b)
			fmt.Printf("CURSOR list: %#v\n", b.itemList)
			if cursor > len(b.itemList)-1 {
				break
			}
			if cursor == -1 {
				break
			}
			var item interface{}
			item = b.itemList[cursor]
			fmt.Printf("BLOCKINGQUEUE: %#v\n", item)
			itemc <- Item{cursor, item}
		}
	}()
	return itemc
}

func (b *BlockingQueue) ItemList() []interface{} {
	b.Lock()
	defer b.Unlock()
	return b.itemList
}

func (b *BlockingQueue) NextLockless() <-chan Item {
	itemc := make(chan Item)
	go func() {
		defer close(itemc)
		var cursor int
		for cursor < len(b.itemList) {
			cursor = <-b.Cursor
			fmt.Println("CURSOR", cursor, len(b.itemList))
			fmt.Printf("CURSOR queue: %#v\n", b)
			fmt.Printf("CURSOR list: %#v\n", b.itemList)
			if cursor > len(b.itemList)-1 {
				break
			}
			if cursor == -1 {
				break
			}
			var item interface{}
			item = b.itemList[cursor]
			fmt.Printf("BLOCKINGQUEUE: %#v\n", item)
			itemc <- Item{cursor, item}
		}
	}()
	return itemc
}

func (b *BlockingQueue) Lock() {
	b.lock.Lock()
}

func (b *BlockingQueue) Unlock() {
	b.lock.Unlock()
}
