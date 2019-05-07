package canalystii

import (
	"fmt"
	"posam/util/concurrentmap"
	"strings"
	"testing"
)

func TestAddInstance(t *testing.T) {
	deviceMap = concurrentmap.NewConcurrentMap()
	d := &Dao{_id: "id"}
	fmt.Println(d)
	addInstance(d)
	if _, found := deviceMap.Get("id"); !found {
		t.Errorf(
			"\nEXPECT: %v\n GET: %v\n\n",
			"instance found",
			"not found",
		)
	}
}

func TestSetID(t *testing.T) {
	cases := []struct {
		id     string
		errmsg string
	}{
		{
			"id",
			"",
		},
		{
			"id",
			"is duplicated",
		},
	}
	deviceMap = concurrentmap.NewConcurrentMap()
	d := &Dao{}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			err := d.setID(c.id)
			if c.errmsg != "" {
				if err == nil {
					t.Errorf("expect error: %v\n", c.errmsg)
					return
				}
				if !strings.Contains(err.Error(), c.errmsg) {
					t.Errorf(
						"\nEXPECT: %v\n GET: %v\n\n",
						c.errmsg,
						err.Error(),
					)
				}
			}
			expect := c.id
			actual := d.id()
			if actual != expect {
				t.Errorf(
					"\nEXPECT: %v\n GET: %v\n\n",
					expect,
					actual,
				)
			}
		})
	}
}
