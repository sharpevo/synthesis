package blockingqueue_test

import (
	"fmt"
	"synthesis/util/blockingqueue"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestBlocking(t *testing.T) {
	testList := []struct {
		v interface{}
		c bool
	}{
		{
			v: 1,
			c: false,
		},
		{
			v: "2",
			c: false,
		},
		{
			v: 0x3,
			c: false,
		},
		{
			v: 4,
			c: true,
		},
	}
	q := blockingqueue.NewBlockingQueue()
	for _, test := range testList {
		if test.c {
			go func() {
				time.Sleep(2 * time.Second)
				fmt.Println("pushing again")
				q.Push(test.v)
				fmt.Println("pushed")
			}()
		} else {
			q.Push(test.v)
		}
	}

	for item := range q.Iter() {
		fmt.Println(item)
	}

	go func() {
		time.Sleep(4 * time.Second)
		q.Reset()
		fmt.Println("reset")
	}()

	for i := range [10]int{} {
		actual, err := q.Pop()
		if err != nil {
			if strings.Contains(err.Error(), "terminated") {
				fmt.Println("error occured as expected")
				return
			} else {
				panic(err)
			}
		}
		test := testList[i]
		if !reflect.DeepEqual(test.v, actual) {
			t.Errorf(
				"\nEXPECT: %q\nGET: %q\n",
				test.v,
				actual,
			)
		}
	}
}

func TestNext(t *testing.T) {
	q := blockingqueue.NewBlockingQueue()
	for i := range [10]int{} {
		v := i
		go func(value int) {
			q.Push(value)
		}(v)
	}
	fmt.Printf("%#v\n", q)
	var wg sync.WaitGroup
	wg.Add(100)
	for i := range [100]int{} {
		count := i
		go func(j int) {
			q.Lock() // cursor may corruption by other goroutines.
			defer q.Unlock()
			go func() {
				q.Cursor <- 0
			}()
			index := 0
			for itemi := range q.NextLockless() {
				fmt.Println(">>>", j, index, itemi)
				index++
				//go func() {
				q.Cursor <- index
				//}()
			}
			wg.Done()
		}(count)
	}
	wg.Wait()
}

func xTestNext(t *testing.T) {
	q := blockingqueue.NewBlockingQueue()
	for i := range [10]int{} {
		v := i
		go func(value int) {
			q.Push(value)
		}(v)
	}
	fmt.Printf("%#v\n", q)
	var wg sync.WaitGroup
	wg.Add(100)
	for i := range [100]int{} {
		count := i
		go func(j int) {
			go func() {
				q.Cursor <- 0
			}()
			index := 0
			for itemi := range q.Next() {
				fmt.Println(">>>", j, index, itemi)
				index++
				//go func() {
				q.Cursor <- index
				//}()
			}
			wg.Done()
		}(count)
	}
	wg.Wait()
}
