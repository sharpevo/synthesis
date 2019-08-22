package concurrentmap_test

import (
	"synthesis/pkg/concurrentmap"
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

func TestDeleteMapEntry(t *testing.T) {
	cmap := concurrentmap.NewConcurrentMap()
	cmap.Set("a", 1)
	cmap.Set("b", "2")
	if err := cmap.Del("a"); err != nil {
		t.Errorf(err.Error())
	}
	for item := range cmap.Iter() {
		t.Log(item)
	}
	if err := cmap.Del("c"); err == nil {
		t.Errorf("error expected")
	}
}

func TestRpl(t *testing.T) {
	cmap := concurrentmap.NewConcurrentMap()
	cmap.Set("key1", 1)
	cmap.Set("key2", false)
	cmap.Set("key3", "test")
	cmap.Set("key4", "test")
	key, err := cmap.Replace(1, 2)
	if err != nil {
		t.Fatal(err)
	}
	if key != "key1" {
		t.Errorf(
			"\nEXPECT: %v\n GET: %v\n\n",
			"key1",
			key,
		)
	}
	key, err = cmap.Replace(false, true)
	if err != nil {
		t.Fatal(err)
	}
	if key != "key2" {
		t.Errorf(
			"\nEXPECT: %v\n GET: %v\n\n",
			"key2",
			key,
		)
	}
	key, err = cmap.Replace("test", "newvalue")
	if err != nil {
		t.Fatal(err)
	}
	if key != "key3" && key != "key4" {
		t.Errorf(
			"\nEXPECT: %v\n GET: %v\n\n",
			"key3/key4",
			key,
		)
	}
	t.Log(cmap)
}
