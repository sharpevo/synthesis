package main

import (
	"posam/dao/canalystii"
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
) (err error) {
	dao, err = canalystii.NewDao(
		devType, devIndex, devID, canIndex, accCode,
		accMask, filter, timing0, timing1, mode,
	)
	return err
}

//export MoveAbsolute
func MoveAbsolute(
	motorCode int,
	position int,
) (resp interface{}, err error) {
	return dao.MoveAbsolute(
		motorCode, position,
	)
}

//export MoveRelative
func MoveRelative(
	motorCode int,
	direction int,
	speed int,
	position int,
) (resp interface{}, err error) {
	return dao.MoveRelative(
		motorCode, direction, speed, position,
	)
}

//export ControlSwitcher
func ControlSwitcher(
	data int,
) (resp interface{}, err error) {
	return dao.ControlSwitcher(data)
}

//export ControlSwitcherAdvanced
func ControlSwitcherAdvanced(
	data int,
	speed int,
	count int,
) (resp interface{}, err error) {
	return dao.ControlSwitcherAdvanced(
		data, speed, count,
	)
}

//export ReadHumiture
func ReadHumiture() (
	resp interface{}, err error,
) {
	return dao.ReadHumiture()
}

//export ReadOxygenConc
func ReadOxygenConc() (
	resp interface{}, err error,
) {
	return dao.ReadOxygenConc()
}

//export ReadPressure
func ReadPressure(
	device int,
) (
	resp interface{}, err error,
) {
	return dao.ReadPressure(device)
}

//export WriteSystemRom
func WriteSystemRom(
	address int,
	value int,
) (resp interface{}, err error) {
	return dao.WriteSystemRom(
		address, value,
	)
}

//export ReadSystemRom
func ReadSystemRom(
	address int,
) (resp interface{}, err error) {
	return dao.ReadSystemRom(address)
}
