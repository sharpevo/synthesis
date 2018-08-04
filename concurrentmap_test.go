package concurrentmap_test

import (
	"posam/util/concurrentmap"
	"testing"
	"time"
)

func TestConcurrencyReadWrite(t *testing.T) {
	cmap := concurrentmap.NewConcurrentMap()
	go func() {
		for {
			cmap.Get("02")
		}
	}()
	go func() {
		for {
			cmap.Set("01", 5)
			cmap.Set("02", "test")
		}
	}()
	select {
	case <-time.After(3 * time.Second):
	}
}
