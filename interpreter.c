#include "_cgo_export.h"

int bridgeCallback(callbackFunc f, struct Instruction* i){
	return f(i);
}


