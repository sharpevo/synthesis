package modbus_test

import (
	//"fmt"
	"posam/protocol/modbus"
	"reflect"
	"testing"
)

func TestCheckSum(t *testing.T) {
	dataList := [][]byte{
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
	crcList := [][]byte{
		[]byte{
			0xD5,
			0xCA,
		},
		[]byte{
			0xE8,
			0x0A,
		},
	}

	for i, v := range dataList {
		crc := modbus.CRC(v)
		if !reflect.DeepEqual(crc, crcList[i]) {
			t.Errorf("\n%d#\nEXPECT: %x\nGET: %x\n", i, crcList[i], crc)
		}
	}

}
