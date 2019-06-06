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
typedef int (__stdcall *f_init_canalyst)(char*, char*, char*, char*, char*, char*,
    char*, char*, char*, char*);
typedef int (__stdcall *f_upsert_variable)(char*, char*);
typedef int (__stdcall *f_get_error)(char*, int);

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
    char msg[30];
    HINSTANCE interpreterlib = LoadLibrary("interpreter.dll");
    if (!interpreterlib) {
        cout << "failed to load dll" << endl;
        return 1;
    }
    f_get_error getLastErrorMessage = (f_get_error)GetProcAddress(
            interpreterlib, "GetLastErrorMessage");
    if (!getLastErrorMessage) {
        cout << "failed to load getLastErrorMessage" << endl;
        return 1;
    }
    f_execute execute = (f_execute)GetProcAddress(interpreterlib, "Execute");
    if (!execute) {
        cout << "failed to load execute" << endl;
        return 1;
    }
    f_register registerHandler = (f_register)GetProcAddress(
            interpreterlib, "RegisterHandlerForInstruction");
    if (!registerHandler) {
        cout << "failed to load register" << endl;
        return 1;
    }
    registerHandler(handler);
    f_init_canalyst initCanalyst = (f_init_canalyst)GetProcAddress(
            interpreterlib, "InitCanalyst");
    if (!initCanalyst) {
        cout << "failed to load initCanalyst" << endl;
        return 1;
    }
    f_upsert_variable upsertVariable = (f_upsert_variable)GetProcAddress(
            interpreterlib, "UpsertVariable");
    if (!upsertVariable) {
        cout << "failed to load upsertVariable" << endl;
        return 1;
    }
    if (!initCanalyst(pchar("4"),pchar("0"),pchar("0x00000001"),
        pchar("0"), pchar("0x00000000"), pchar("0xFFFFFFFF"), pchar("0"),
        pchar("0x00"), pchar("0x1c"), pchar("0"))){
        getLastErrorMessage(&msg[0], 30);
        cout << "C++: get error from DLL when init canalyst: " << msg << endl;
        return 1;
    }
    upsertVariable(pchar("WB_NPV"), pchar("4097"));
    char* switch_args[3] = {pchar("FRAME_ID"), pchar("val"), pchar("WB_NPV")};
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
    char* unknown_args[2] = {pchar("ARG1"), pchar("ARG2")};
    instruction i3 = {
        .name = pchar("FOO"),
        .arguments = unknown_args,
        .argumentCount = 2,
        .ignoreError = 0,
        .remark = pchar("instruction-3"),
    };
    int return3;
    thread t3{[&i3, &return3, execute]{
        return3 = execute(&i3);
    }};
    t0.join();
    t1.join();
    t2.join();
    t3.join();
    if (!return3){
        char msg3[30];
        getLastErrorMessage(&msg3[0], 30);
        cout << "C++: get error from DLL when execute instruction-3: " << msg3 << endl;
    }
}
