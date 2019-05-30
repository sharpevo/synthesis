package main

/*
#include <stdlib.h>
typedef int(*callbackFunc)(char*, char*);
struct Instruction{
	char* Name;
	char** Arguments;
	int IgnoreError;
	callbackFunc callback;
};
extern int bridgeCallback(callbackFunc f, char* a, char* b);
*/
import "C"

import (
	"log"
	//"unsafe"
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
	//defer C.free(unsafe.Pointer(resp))
	//defer C.free(unsafe.Pointer(err))
	f := C.callbackFunc(instruction.callback)
	rst := int(C.bridgeCallback(f, resp, err))
	log.Println("callback: ", rst)
}
