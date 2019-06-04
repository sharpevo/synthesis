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
	"log"
	"posam/dao"
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

var stack = instruction.NewStack()

func main() {}

//export NewInstruction
func NewInstruction(instruction *C.struct_instruction) int {
	return 0
}

//export Execute
func Execute(i *C.struct_instruction) {
	//func Execute(group *C.struct_instruction_set) {

	if instance, _ := canalystii.Instance("1"); instance != nil {
		//log.Printf(">>>> Device %q has been initialized\n", "1")
		// TODO: register device tree
	} else {
		//log.Printf(">>>> use previous")
		// CAN
		devtype := "4"
		devindex := "0"
		frameid := "0x00000001"
		canindex := "0"
		acccode := "0x00000000"
		accmask := "0xFFFFFFFF"
		filter := "0"
		timing0 := "0x00"
		timing1 := "0x1C"
		mode := "0"

		devtypev, _ := vrb.NewVariable(canalystii.DEVICE_TYPE, "4")
		devindexv, _ := vrb.NewVariable(canalystii.DEVICE_INDEX, "0")
		frameidv, _ := vrb.NewVariable(canalystii.FRAME_ID, "0x00000001")
		canindexv, _ := vrb.NewVariable(canalystii.CAN_INDEX, "0")
		acccodev, _ := vrb.NewVariable(canalystii.ACC_CODE, "0x00000000")
		accmaskv, _ := vrb.NewVariable(canalystii.ACC_MASK, "0xFFFFFFFF")
		filterv, _ := vrb.NewVariable(canalystii.FILTER, "0")
		timing0v, _ := vrb.NewVariable(canalystii.TIMING0, "0x00")
		timing1v, _ := vrb.NewVariable(canalystii.TIMING1, "0x1C")
		modev, _ := vrb.NewVariable(canalystii.MODE, "0")

		stack.Set(devtypev)
		stack.Set(devindexv)
		stack.Set(frameidv)
		stack.Set(canindexv)
		stack.Set(acccodev)
		stack.Set(accmaskv)
		stack.Set(filterv)
		stack.Set(timing0v)
		stack.Set(timing1v)
		stack.Set(modev)

		terminatecc := make(chan chan interface{}, 1)
		defer close(terminatecc)

		instructionMap := dao.NewInstructionMap()
		for item := range canalystii.InstructionMap.Iter() {
			instructionMap.Set(item.Key, item.Value.(reflect.Type))
		}
		for item := range dao.InstructionMap.Iter() {
			instructionMap.Set(item.Key, item.Value.(reflect.Type))
		}
		interpreter.InitParser(instructionMap)
		if _, err := canalystii.NewDao(
			devtype,
			devindex,
			frameid,
			canindex,
			acccode,
			accmask,
			filter,
			timing0,
			timing1,
			mode,
		); err != nil {
			log.Fatal(">>>>", err)
		}
	}

	// parse

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
		rst := int(C.bridgeCallback(handlerForInstruction, i))
		log.Println("MSG: ", resp.Output, resp.Error, rst)
	}
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
