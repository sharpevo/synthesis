package ricoh_g5_test

import (
	"reflect"
	"synthesis/internal/dao"
	"synthesis/internal/dao/ricoh_g5"
	"testing"
)

func TestPdu(t *testing.T) {
	unitList := []dao.Unit{
		ricoh_g5.ErrorCodeUnit,
		ricoh_g5.PrinterStatusUnit,
		ricoh_g5.PrintDataUnit,
		ricoh_g5.WaveformUnit,
	}

	expectedList := [][]byte{
		[]byte{
			0x01,
			0x00,
			0x00,
			0x00,
		},
		[]byte{
			0x02,
			0x00,
			0x00,
			0x00,
		},
		[]byte{
			0x03,
			0x00,
			0x00,
			0x00,
		},
		[]byte{
			0x04,
			0x00,
			0x00,
			0x00,
		},
	}
	for i, v := range unitList {
		req := v.Request()
		if !reflect.DeepEqual(req.Bytes(), expectedList[i]) {
			t.Errorf("\n%d#\nEXPECT: %x\nGET: %x\n", i, expectedList[i], req.Bytes())
		}
	}
}
