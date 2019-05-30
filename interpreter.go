package main

/*
#include <stdlib.h>
struct instruction{
	char* name;
	char** arguments;
	int ignoreError;
	char* output;
	char* err;
};
typedef int(*f_handlerInstruction)(struct instruction*);
extern int bridgeCallback(f_handlerInstruction f, struct instruction* i);
*/
import "C"

import (
	"log"
	"unsafe"
	//"posam/interpreter"
)

func main() {}

//export NewInstruction
func NewInstruction(instruction *C.struct_instruction) int {
	return 0
}

//export Execute
func Execute(instruction *C.struct_instruction) {
	log.Println("executing", C.GoString(instruction.name))
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
