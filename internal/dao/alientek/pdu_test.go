package alientek_test

import (
	"reflect"
	"synthesis/internal/dao"
	"synthesis/internal/dao/alientek"
	"testing"
)

func TestPdu(t *testing.T) {
	unitList := []dao.Unit{
		alientek.LedOnUnit,
		alientek.LedOffUnit,
	}
	expectList := [][]byte{
		[]byte{
			0x01,
			0x03,
			0x00,
			0x01,
			0x00,
			0x01,
		},
		[]byte{
			0x01,
			0x02,
			0x00,
			0x01,
			0x00,
			0x01,
		},
	}

	for i, v := range unitList {
		req := v.Request()
		req.Address = 0x01
		if !reflect.DeepEqual(req.Bytes(), expectList[i]) {
			t.Errorf("\n%d#\nEXPECT: %x\nGET: %x\n", i, expectList[i], req.Bytes())
		}
	}
}
