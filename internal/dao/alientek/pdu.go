package alientek

import (
	"synthesis/internal/dao"
)

var LedOnUnit,
	LedOffUnit dao.Unit

func init() {

	// LedOnUnit{{{
	LedOnRequest := dao.Request{}
	LedOnRequest.Function = 0x03
	LedOnRequest.Arguments = []byte{
		0x00,
		0x01,
		0x00,
		0x01,
	}
	LedOnRecResp := []byte{0x55}
	LedOnComResp := []byte{
		0x01,
		0x83,
		0x02,
		0xC0,
		0xF1,
	}

	LedOnUnit.SetRequest(LedOnRequest)
	LedOnUnit.SetRecResp(LedOnRecResp)
	LedOnUnit.SetComResp(LedOnComResp)
	// }}}

	// LedOffUnit{{{
	LedOffRequest := dao.Request{}
	LedOffRequest.Function = 0x02
	LedOffRequest.Arguments = []byte{
		0x00,
		0x01,
		0x00,
		0x01,
	}
	LedOffRecResp := []byte{0x55}
	LedOffComResp := []byte{
		0x01,
		0x82,
		0x02,
		0xC1,
		0x61,
	}
	LedOffUnit.SetRequest(LedOffRequest)
	LedOffUnit.SetRecResp(LedOffRecResp)
	LedOffUnit.SetComResp(LedOffComResp)
	// }}}

}
