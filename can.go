package main

import (
	"fmt"
	"posam/dao/canalystii"
	"reflect"
)
import "C"

func main() {}

var dao *canalystii.Dao

//export NewDao
func NewDao(
	devType string,
	devIndex string,
	devID string,
	canIndex string,
	accCode string,
	accMask string,
	filter string,
	timing0 string,
	timing1 string,
	mode string,
) int {
	var err error
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
func ReadHumiture() []float64 {
	resp, err := dao.ReadHumiture()
	if err != nil {
		fmt.Println(err)
		return []float64{0}
		// TODO: error handling
	}
	fmt.Println("ReadHumiture:", resp)
	output := resp.([]float64)
	return output
}

func ReadHumiture2() uintptr {
	//func ReadHumiture() []float64 {
	resp, err := dao.ReadHumiture()
	if err != nil {
		fmt.Println(err)
		return []float64{0}
		// TODO: error handling
	}
	fmt.Println("ReadHumiture:", resp)
	output := resp.([]float64)

	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&output))
	return hdr.Data
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
