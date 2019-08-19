package main

/*
#include <stdlib.h>
struct instruction{
	char* name;
	char** arguments;
	int argumentCount;
	int ignoreError;
	char* output;
	char* err;
	char* remark;

	int executionType;
	int instructionCount;
	struct instruction** instructions;
};
typedef int(*f_handlerInstruction)(struct instruction*);
extern int bridgeCallback(f_handlerInstruction f, struct instruction* i);
*/
import "C"

import (
	"fmt"
	//"io/ioutil"
	"log"
	"posam/dao"
	"posam/dao/aoztech"
	"posam/dao/canalystii"
	"posam/instruction"
	"posam/interpreter"
	"posam/interpreter/vrb"
	//"posam/util"
	"sync"
	//"posam/util/blockingqueue"
	"reflect"
	"unsafe"
)

const (
	SERIAL      = 0
	CONCURRENCY = 1
)

var errmsg string

var stack = instruction.NewStack()
var lock = sync.Mutex{}
var instructionMap = dao.NewInstructionMap()
var launched = false

func init() {
	//log.SetOutput(ioutil.Discard)
}
func main() {}

//export NewInstruction
func NewInstruction(instruction *C.struct_instruction) int {
	return 0
}

func LaunchInterpreter() {
	lock.Lock()
	defer lock.Unlock()
	if !launched {
		for item := range dao.InstructionMap.Iter() {
			instructionMap.Set(item.Key, item.Value.(reflect.Type))
		}
		interpreter.InitParser(instructionMap)
		launched = true
	}
}

//export UpsertVariable
func UpsertVariable(name *C.char, value *C.char) int {
	return upsertVariable(C.GoString(name), C.GoString(value))
}

func upsertVariable(name string, value string) int {
	v, err := vrb.NewVariable(name, value)
	if err != nil {
		errmsg = err.Error()
		log.Println(err)
		return 0
	}
	stack.Set(v)
	return 1
}

//export InitCanalyst
func InitCanalyst(
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
	for item := range canalystii.InstructionMap.Iter() {
		instructionMap.Set(item.Key, item.Value.(reflect.Type))
	}

	if _, err := canalystii.NewDao(
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
	); err != nil {
		errmsg = err.Error()
		log.Println(err)
		return 0
	}
	return 1
}

//export InitAoztech
func InitAoztech(
	nameChar *C.char,
	baudChar *C.char,
	axisXIDChar *C.char,
	axisXSetupFileChar *C.char,
	axisYIDChar *C.char,
	axisYSetupFileChar *C.char,
) int {
	name,
		baud,
		axisXID,
		axisXSetupFile,
		axisYID,
		axisYSetupFile := C.GoString(nameChar),
		C.GoString(baudChar),
		C.GoString(axisXIDChar),
		C.GoString(axisXSetupFileChar),
		C.GoString(axisYIDChar),
		C.GoString(axisYSetupFileChar)
	for item := range aoztech.InstructionMap.Iter() {
		instructionMap.Set(item.Key, item.Value.(reflect.Type))
	}
	if _, err := aoztech.NewDao(
		name,
		baud,
		axisXID,
		axisXSetupFile,
		axisYID,
		axisYSetupFile,
	); err != nil {
		errmsg = err.Error()
		log.Println(err)
		return 0
	}
	return 1
}

//export Execute
func Execute(i *C.struct_instruction) int {
	LaunchInterpreter()
	//fmt.Println(">> count", int(i.instructionCount))
	if int(i.instructionCount) > 0 {
		//fmt.Println(">> group")
		count := int(i.instructionCount)
		itemsC := (*[1 << 30]*C.struct_instruction)(unsafe.Pointer(i.instructions))[:count:count]
		switch int(i.executionType) {
		case SERIAL:
			for _, item := range itemsC {
				Execute(item)
			}
		case CONCURRENCY:
			var wg sync.WaitGroup
			for _, item := range itemsC {
				wg.Add(1)
				go func(w *sync.WaitGroup, itm *C.struct_instruction) {
					Execute(itm)
					w.Done()
				}(&wg, item)
			}
			wg.Wait()
		}
	} else {
		//fmt.Println(">> single", C.GoString(i.name))
		count := int(i.argumentCount)
		//log.Println("count", count)
		argsC := (*[1 << 30]*C.char)(unsafe.Pointer(i.arguments))[:count:count]
		//log.Println("argsC", argsC)
		args := make([]string, count)
		for i, s := range argsC {
			//log.Println("argC", s)
			args[i] = C.GoString(s)
			//log.Println("arg", args[i])
		}

		ignoreErrorI := int(i.ignoreError)
		ignoreError := ignoreErrorI == 0
		statement := interpreter.Statement{
			Stack:           stack,
			InstructionName: C.GoString(i.name),
			Arguments:       args,
			IgnoreError:     ignoreError,
		}

		//util.Go(func() {
		terminatec := make(chan interface{})
		completec := make(chan interface{})
		//util.Go(func() {
		//<-completec
		//})
		resp := <-statement.Execute(terminatec, completec)
		output := C.CString(fmt.Sprintf("%v", resp.Output))
		err := C.CString(fmt.Sprintf("%v", resp.Error))
		defer C.free(unsafe.Pointer(output))
		defer C.free(unsafe.Pointer(err))
		i.output = output
		i.err = err
		//if resp.Error != nil {
		//errmsg = resp.Error.Error()
		//}
		rst := int(C.bridgeCallback(handlerForInstruction, i))
		if rst == 0 {
			// TODO: error occurs at handlerForInstruction
		}
		//log.Println("MSG: ", resp.Output, resp.Error, rst)

		// only for testing
		if C.GoString(i.name) == "FOO" {
			errmsg = "BAR"
			return 0
		}
	}
	return 1
}

func ExecuteInstruction(instruction *C.struct_instruction) {
	log.Println("executing", C.GoString(instruction.name))

	count := int(instruction.argumentCount)
	log.Println("count", count)
	argsC := (*[1 << 30]*C.char)(unsafe.Pointer(instruction.arguments))[:count:count]
	log.Println("argsC", argsC)
	args := make([]string, count)
	for i, s := range argsC {
		log.Println("argC", s)
		args[i] = C.GoString(s)
		log.Println("arg", args[i])
	}
	log.Printf("args: %#v\n", args)

	resp := C.CString("moved to 10mm")
	err := C.CString("no error")
	defer C.free(unsafe.Pointer(resp))
	defer C.free(unsafe.Pointer(err))
	instruction.output = resp
	instruction.err = err
	rst := int(C.bridgeCallback(handlerForInstruction, instruction))
	log.Println("handler: ", rst)
}

var handlerForInstruction C.f_handlerInstruction

//export RegisterHandlerForInstruction
func RegisterHandlerForInstruction(h C.f_handlerInstruction) {
	handlerForInstruction = h
}

//export GetLastErrorMessage
func GetLastErrorMessage(p *C.char, n int) int {
	if n < 0 {
		return 0
	}
	if n > len(errmsg)-1 {
		n = len(errmsg)
	}
	pp := (*[1 << 30]byte)(unsafe.Pointer(p))
	copy(pp[:], errmsg[:n])
	pp[n] = 0
	return 1
}
