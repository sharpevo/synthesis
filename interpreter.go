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
};
typedef int(*f_handlerInstruction)(struct instruction*);
extern int bridgeCallback(f_handlerInstruction f, struct instruction* i);
*/
import "C"

import (
	"log"
	//"posam/interpreter"
	"unsafe"
)

func main() {}

//export NewInstruction
func NewInstruction(instruction *C.struct_instruction) int {
	return 0
}

//export Execute
func Execute(instruction *C.struct_instruction) {
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
