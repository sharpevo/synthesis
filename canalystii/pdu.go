package canalystii

import (
	"posam/dao"
)

var SensorOxygenConcUnit,
	SensorHumitureUnit dao.Unit

func init() {
	// SensorOxygenConcUnit{{{
	SensorOxygenConcRequest := dao.Request{}
	SensorOxygenConcRequest.Function = 0x0D
	SensorOxygenConcRequest.Arguments = []byte{
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
	}
	SensorOxygenConcUnit.SetRequest(SensorOxygenConcRequest)
	SensorOxygenConcUnit.SetRecResp([]byte{})
	SensorOxygenConcUnit.SetComResp([]byte{})
	// }}}
}
