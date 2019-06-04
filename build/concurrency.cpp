#include <windows.h>
#include <iostream>
#include <thread>
extern "C"{
#include "interpreter.h"
}
using namespace std;
typedef int (__stdcall *f_handler)(instruction*);
typedef int (__stdcall *f_register)(f_handler);
typedef int (__stdcall *f_execute)(instruction*);

const int SERIAL = 0;
const int CONCURRENCY = 1;

char* pchar(string input){
    return const_cast<char*>(input.c_str());
}

int handler(instruction* i){
    cout << "==== handler in C++ ====" << endl;
    cout << "instruction '" << i->name << "' is completed" << endl;
    cout << "output: " << i->output << endl;
    cout << "error: " << i->err << endl;
    cout << "remark: " << i->remark << endl;
    cout << "========================" << endl;
    return 0;
}

int main(){
    HINSTANCE interpreterlib = LoadLibrary("interpreter.dll");
    if (!interpreterlib) {
        cout << "failed to load dll" << endl;
        return 1;
    }
    f_execute execute = (f_execute)GetProcAddress(interpreterlib, "Execute");
    if (!execute) {
        cout << "failed to load execute" << endl;
        return 1;
    }
    f_register registerHandler = (f_register)GetProcAddress(interpreterlib, "RegisterHandlerForInstruction");
    if (!registerHandler) {
        cout << "failed to load register" << endl;
        return 1;
    }
    registerHandler(handler);
    char* switch_args[3] = {pchar("FRAME_ID"), pchar("val"), pchar("12")};
    instruction i0 = {
        .name = pchar("SWITCH"),
        .arguments = switch_args,
        .argumentCount = 3,
        .ignoreError = 0,
        .remark = pchar("instruction-0"),
    };
    thread t0(execute, &i0);
    char* sleep_args[1] = {pchar("3")};
    instruction i1 = {
        .name = pchar("SLEEP"),
        .arguments = sleep_args,
        .argumentCount = 1,
        .ignoreError = 0,
        .remark = pchar("instruction-1"),
    };
    thread t1(execute, &i1);
    char* humiture_args[3] = {pchar("FRAME_ID"), pchar("humi"), pchar("temp")};
    instruction i2 = {
        .name = pchar("HUMITURE"),
        .arguments = humiture_args,
        .argumentCount = 3,
        .ignoreError = 0,
        .remark = pchar("instruction-2"),
    };
    thread t2(execute, &i2);
    t0.join();
    t1.join();
    t2.join();
}
