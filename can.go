package main

import (
	"fmt"
	"posam/dao/canalystii"
	//"reflect"
)
import "C"

func main() {}

var dao *canalystii.Dao

//export NewDao
func NewDao(
	devTypeChar *C.char,
	devIndexChar *C.char,
	devIDChar *C.char,
	canIndexChar *C.char,
	accCodeChar *C.char,
	accMaskChar *C.char,
	filterChar *C.char,
	timing0Char *C.char,
	timing1Char *C.char,
	modeChar *C.char,
) int {
	devType, devIndex, devID, canIndex, accCode,
		accMask, filter, timing0, timing1, mode := C.GoString(devTypeChar),
		C.GoString(devIndexChar), C.GoString(devIDChar), C.GoString(canIndexChar), C.GoString(accCodeChar), C.GoString(accMaskChar), C.GoString(filterChar), C.GoString(timing0Char), C.GoString(timing1Char), C.GoString(modeChar)
	var err error
	fmt.Println(devType)
	fmt.Println(devID)
	dao, err = canalystii.NewDao(
		devType, devIndex, devID, canIndex, accCode,
		accMask, filter, timing0, timing1, mode,
	)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return 1
}

//export MoveAbsolute
func MoveAbsolute(
	motorCode int,
	position int,
) int {
	resp, err := dao.MoveAbsolute(
		motorCode, position,
	)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	fmt.Println("MoveAbosute:", resp)
	return 1
}

//export MoveRelative
func MoveRelative(
	motorCode int,
	direction int,
	speed int,
	position int,
) int {
	resp, err := dao.MoveRelative(
		motorCode, direction, speed, position,
	)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	fmt.Println("MoveRelative:", resp)
	return 1
}

//export ControlSwitcher
func ControlSwitcher(data int) int {
	resp, err := dao.ControlSwitcher(data)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	fmt.Println("ControSwitcher:", resp)
	return 1
}

//export ControlSwitcherAdvanced
func ControlSwitcherAdvanced(
	data int,
	speed int,
	count int,
) int {
	resp, err := dao.ControlSwitcherAdvanced(
		data, speed, count,
	)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	fmt.Println("ControSwitcher:", resp)
	return 1
}

//export ReadHumiture
func ReadHumiture(temp *float64, humi *float64) int {
	resp, err := dao.ReadHumiture()
	if err != nil {
		fmt.Println(err)
		return 0
	}
	output, ok := resp.([]float64)
	if !ok {
		fmt.Println("Error: invaild humiture data type")
		return 1
	}
	if len(output) != 2 {
		fmt.Println("Error: invaild humiture data")
		return 1
	}
	fmt.Println("ReadHumiture:", output[0], output[1])
	*temp = output[0]
	*humi = output[1]
	return 1
}

//export Test
func Test(inputChar *C.char, output **C.char) int {
	//fmt.Println(C.GoString(Input))
	//*Output = C.CString(fmt.Sprintf("From DLL: Hello, %s!\n", C.GoString(Input)))
	//fmt.Println("Message: ", C.GoString(*Output))
	//return 1

	input := C.GoString(inputChar)
	fmt.Println("GO1", input)
	rst := C.CString("中文" + input)
	*output = rst
	fmt.Println("GO2", C.GoString(*output))
	//c.free(unsafe.Pointer(rst))
	return int(len(C.GoString(*output)))
}

//export ReadOxygenConc
func ReadOxygenConc() float64 {
	resp, err := dao.ReadOxygenConc()
	if err != nil {
		fmt.Println(err)
		return float64(0)
	}
	fmt.Println("ReadOxygenConc:", resp)
	output := resp.(float64)
	return output
}

//export ReadPressure
func ReadPressure(device int) int64 {
	resp, err := dao.ReadPressure(device)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	fmt.Println("ReadPressure:", resp)
	output := resp.(int64)
	return output
}

//export WriteSystemRom
func WriteSystemRom(address int, value int) int {
	resp, err := dao.WriteSystemRom(address, value)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	fmt.Println("WriteSystemRom:", resp)
	return 1
}

//export ReadSystemRom
func ReadSystemRom(address int) int {
	resp, err := dao.ReadSystemRom(address)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	fmt.Println("WriteSystemRom:", resp)
	return 1
}
