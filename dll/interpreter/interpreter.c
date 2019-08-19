#include "_cgo_export.h"

int bridgeCallback(f_handlerInstruction h, struct instruction* i){
	return h(i);
}


