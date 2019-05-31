#include <windows.h>
#include <iostream>
extern "C"{
#include "interpreter.h"
}
typedef int (__stdcall *f_handler)(instruction*);
typedef int (__stdcall *f_register)(f_handler);
typedef int (__stdcall *f_execute)(instruction*);

const int SERIAL = 0;
const int CONCURRENCY = 1;

char* pchar(std::string input){
    return const_cast<char*>(input.c_str());
}

int handler(instruction* i){
    std::cout << "==== handler in C++ ====" << std::endl;
    std::cout << "instruction '" << i->name << "' is completed" << std::endl;
    std::cout << "output: " << i->output << std::endl;
    std::cout << "error: " << i->err << std::endl;
    //std::cout << "id: " << i->custom << std::endl;
    std::cout << "========================" << std::endl;
    return 22;
}

int main(){
    HINSTANCE interpreterlib = LoadLibrary("interpreter.dll");
    if (!interpreterlib) {
        std::cout << "failed to load dll" << std::endl;
        return 1;
    }
    f_execute execute = (f_execute)GetProcAddress(interpreterlib, "Execute");
    if (!execute) {
        std::cout << "failed to load execute" << std::endl;
        return 1;
    }
    f_register registerHandler = (f_register)GetProcAddress(interpreterlib, "RegisterHandlerForInstruction");
    if (!registerHandler) {
        std::cout << "failed to load register" << std::endl;
        return 1;
    }
    registerHandler(handler);

    instruction s = {
        .executionType = SERIAL,
        .instructionCount = 2,
        .instructions = new struct instruction* [2],
    };

    char* name = pchar("SWITCH");
    int argCount = 3;
    char** args = new char* [argCount];
    args[0] = pchar("FRAME_ID");
    args[1] = pchar("val");
    args[2] = pchar("12");

    instruction i0 = {
        .name = name,
        .arguments = args,
        .argumentCount = argCount,
        .ignoreError = 0,
    };
    s.instructions[0] = &i0;
        //.custom = pchar("id-1"),
    //execute(&i);

    name = pchar("HUMITURE");
    argCount = 3;
    args = new char* [argCount];
    args[0] = pchar("FRAME_ID");
    args[1] = pchar("humi");
    args[2] = pchar("temp");

    instruction i1 = {
        .name = name,
        .arguments = args,
        .argumentCount = argCount,
        .ignoreError = 0,
    };
    s.instructions[1] = &i1;
        //.custom = pchar("id-2"),
    //execute(&i);
    execute(&s);
}
