package blockingqueue_test

import (
	"fmt"
	"posam/util/blockingqueue"
	"reflect"
	"strings"
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
				fmt.Println("pushed again")
				q.Push(test.v)
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
