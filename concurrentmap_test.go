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

func TestNewConcurrentMap(t *testing.T) {
	cmap := concurrentmap.NewConcurrentMap()
	cmap.Set("a", 1)
	cmap.Set("b", "2")
	newCmap := concurrentmap.NewConcurrentMap(cmap)
	newCmap.Set("c", 3.0)
	t.Logf("%#v", newCmap)

	for item := range cmap.Iter() {
		if v, _ := newCmap.Get(item.Key); v != item.Value {
			t.Errorf(
				"\nEXPECT: %v\n GET: %v\n\n",
				item.Value,
				v,
			)
		}
	}
}
