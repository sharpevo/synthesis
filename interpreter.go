package main

/*
#include <stdlib.h>
struct Instruction{
	char* Name;
	char** Arguments;
	int IgnoreError;
	char* Output;
	char* Error;
};
typedef int(*callbackFunc)(struct Instruction*);
extern int bridgeCallback(callbackFunc f, struct Instruction* i);
*/
import "C"

import (
	"log"
	"unsafe"
	//"posam/interpreter"
)

func main() {}

//type Instruction struct {
//interpreter.Statement
//}

//type response_callback func(string string) int

//export NewInstruction
func NewInstruction(instruction *C.struct_Instruction) int {
	return 0
}

//export Execute
func Execute(instruction *C.struct_Instruction) {
	log.Println("executing", C.GoString(instruction.Name))
	resp := C.CString("response string")
	err := C.CString("error string")
	defer C.free(unsafe.Pointer(resp))
	defer C.free(unsafe.Pointer(err))
	instruction.Output = resp
	instruction.Error = err
	//f := C.callbackFunc(instruction.callback)
	rst := int(C.bridgeCallback(handlerForInstruction, instruction))
	log.Println("callback: ", rst)
}

var handlerForInstruction C.callbackFunc

//export RegisterHandlerForInstruction
func RegisterHandlerForInstruction(f C.callbackFunc) {
	handlerForInstruction = f
}
