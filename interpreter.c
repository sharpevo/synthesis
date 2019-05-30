#include "_cgo_export.h"

int bridgeCallback(callbackFunc f, char* a, char* b){
	return f(a, b);
}


