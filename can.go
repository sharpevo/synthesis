package main

import "C"

import (
	"log"
	"posam/dao/canalystii"
	//"reflect"
)

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
	devType,
		devIndex,
		devID,
		canIndex,
		accCode,
		accMask,
		filter,
		timing0,
		timing1,
		mode := C.GoString(devTypeChar),
		C.GoString(devIndexChar),
		C.GoString(devIDChar),
		C.GoString(canIndexChar),
		C.GoString(accCodeChar),
		C.GoString(accMaskChar),
		C.GoString(filterChar),
		C.GoString(timing0Char),
		C.GoString(timing1Char),
		C.GoString(modeChar)
	var err error
	//fmt.Println(devType)
	//fmt.Println(devID)
	dao, err = canalystii.NewDao(
		devType,
		devIndex,
		devID,
		canIndex,
		accCode,
		accMask,
		filter,
		timing0,
		timing1,
		mode,
	)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return 0
	}
	log.Println("MoveAbosute:", resp)
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
		log.Println(err)
		return 0
	}
	log.Println("MoveRelative:", resp)
	return 1
}

//export ControlSwitcher
func ControlSwitcher(data int) int {
	resp, err := dao.ControlSwitcher(data)
	if err != nil {
		log.Println(err)
		return 0
	}
	log.Println("ControSwitcher:", resp)
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
		log.Println(err)
		return 0
	}
	log.Println("ControSwitcher:", resp)
	return 1
}

//export ReadHumiture
func ReadHumiture(temp *float64, humi *float64) int {
	resp, err := dao.ReadHumiture()
	if err != nil {
		log.Println(err)
		return 0
	}
	output, ok := resp.([]float64)
	if !ok {
		log.Println("Error: invaild humiture data type")
		return 1
	}
	if len(output) != 2 {
		log.Println("Error: invaild humiture data")
		return 1
	}
	log.Println("ReadHumiture:", output[0], output[1])
	*temp = output[0]
	*humi = output[1]
	return 1
}

//func Test(inputChar *C.char, output *string) int {
////fmt.Println(C.GoString(Input))
///[>Output = C.CString(fmt.Sprintf("From DLL: Hello, %s!\n", C.GoString(Input)))
////fmt.Println("Message: ", C.GoString(*Output))
////return 1

//input := C.GoString(inputChar)
//fmt.Println("GO1", input)
//*output = "中文" + input
//fmt.Println("GO2", *output)
//return 1
//}

//export Test
func Test(inputChar *C.char, output **C.char) int {
	//fmt.Println(C.GoString(Input))
	//*Output = C.CString(fmt.Sprintf("From DLL: Hello, %s!\n", C.GoString(Input)))
	//fmt.Println("Message: ", C.GoString(*Output))
	//return 1

	input := C.GoString(inputChar)
	log.Println("GO1", input)
	rst := C.CString("中文" + input)
	*output = rst
	log.Println("GO2", C.GoString(*output))
	//C.free(unsafe.Pointer(rst))
	return int(len(C.GoString(*output)))
}

//export ReadOxygenConc
func ReadOxygenConc(output *float64) int {
	resp, err := dao.ReadOxygenConc()
	if err != nil {
		log.Println(err)
		return 0
	}
	log.Println("ReadOxygenConc:", resp)
	*output = resp.(float64)
	return 1
}

//export ReadPressure
func ReadPressure(device int, output *int64) int {
	resp, err := dao.ReadPressure(device)
	if err != nil {
		log.Println(err)
		return 0
	}
	log.Println("ReadPressure:", resp)
	*output = resp.(int64)
	return 1
}

//export WriteSystemRom
func WriteSystemRom(address int, value int) int {
	resp, err := dao.WriteSystemRom(address, value)
	if err != nil {
		log.Println(err)
		return 0
	}
	log.Println("WriteSystemRom:", resp)
	return 1
}

//export ReadSystemRom
func ReadSystemRom(address int, output *int64) int {
	resp, err := dao.ReadSystemRom(address)
	if err != nil {
		log.Println(err)
		return 0
	}
	log.Println("WriteSystemRom:", resp)
	*output = resp.(int64)
	return 1
}
