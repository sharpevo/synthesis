package canalystii

import (
	"synthesis/internal/dao"
	"synthesis/internal/protocol/usbcan"
)

// The following variables defines basic protocol data units according to the
// documents.
var MotorMoveRelativeUnit,
	MotorMoveAbsoluteUnit,
	MotorResetUnit,
	SwitcherControlUnit,
	SwitcherControlAdvancedUnit,
	SensorHumitureUnit,
	SensorOxygenConcUnit,
	SensorPressureUnit,
	SystemRomReadUnit,
	SystemRomWriteUnit dao.Unit

func init() {
	// MotorMoveRelativeUnit{{{
	MotorMoveRelativeRequest := dao.Request{}
	MotorMoveRelativeRequest.Function = 0x01
	MotorMoveRelativeUnit.SetRequest(MotorMoveRelativeRequest)
	MotorMoveRelativeUnit.SetRecResp(responseReceived())
	MotorMoveRelativeUnit.SetComResp(responseCompleted())
	// }}}

	// MotorResetUnit{{{
	MotorResetRequest := dao.Request{}
	MotorResetRequest.Function = 0x02
	MotorResetUnit.SetRequest(MotorResetRequest)
	MotorResetUnit.SetRecResp(responseReceived())
	MotorResetUnit.SetComResp(responseCompleted())
	// }}}

	// MotorMoveAbsoluteUnit{{{
	MotorMoveAbsoluteRequest := dao.Request{}
	MotorMoveAbsoluteRequest.Function = 0x03
	MotorMoveAbsoluteUnit.SetRequest(MotorMoveAbsoluteRequest)
	MotorMoveAbsoluteUnit.SetRecResp(responseReceived())
	MotorMoveAbsoluteUnit.SetComResp(responseCompleted())
	// }}}

	// SwitcherControlUnit{{{
	SwitcherControlRequest := dao.Request{}
	SwitcherControlRequest.Function = 0x0A
	SwitcherControlUnit.SetRequest(SwitcherControlRequest)
	SwitcherControlUnit.SetRecResp(responseNil())
	SwitcherControlUnit.SetComResp(responseNil())
	// }}}

	// SwitcherControlAdvancedUnit{{{
	SwitcherControlAdvancedRequest := dao.Request{}
	SwitcherControlAdvancedRequest.Function = 0x0B
	SwitcherControlAdvancedUnit.SetRequest(SwitcherControlAdvancedRequest)
	SwitcherControlAdvancedUnit.SetRecResp(responseReceived())
	SwitcherControlAdvancedUnit.SetComResp(responseCompleted())
	// }}}

	// SensorHumitureUnit{{{
	SensorHumitureRequest := dao.Request{}
	SensorHumitureRequest.Function = 0x0C
	SensorHumitureRequest.Arguments = []byte{
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
	}
	SensorHumitureUnit.SetRequest(SensorHumitureRequest)
	SensorHumitureUnit.SetRecResp(responseNil())
	SensorHumitureUnit.SetComResp(responseNil())
	// }}}

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
	SensorOxygenConcUnit.SetRecResp(responseNil())
	SensorOxygenConcUnit.SetComResp(responseNil())
	// }}}

	// SensorPressure{{{
	SensorPressureRequest := dao.Request{}
	SensorPressureRequest.Function = 0x0E
	SensorPressureUnit.SetRequest(SensorPressureRequest)
	SensorPressureUnit.SetRecResp(responseNil())
	SensorPressureUnit.SetComResp(responseNil())
	// }}}

	// SystemRomReadUnit{{{
	SystemRomReadRequest := dao.Request{}
	SystemRomReadRequest.Function = 0xF0
	SystemRomReadUnit.SetRequest(SystemRomReadRequest)
	SystemRomReadUnit.SetRecResp(responseNil())
	SystemRomReadUnit.SetComResp(responseCompleted())
	// }}}

	// SystemRomWriteUnit{{{
	SystemRomWriteRequest := dao.Request{}
	SystemRomWriteRequest.Function = 0xF1
	SystemRomWriteUnit.SetRequest(SystemRomReadRequest)
	SystemRomWriteUnit.SetRecResp(responseNil())
	SystemRomWriteUnit.SetComResp(responseCompleted())
	// }}}
}

func responseReceived() []byte {
	return []byte{usbcan.STATUS_CODE_RECEIVED}
}

func responseCompleted() []byte {
	return []byte{usbcan.STATUS_CODE_COMPLETED}
}

func responseNil() []byte {
	return []byte{}
}
