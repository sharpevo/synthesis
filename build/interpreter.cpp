#include <windows.h>
#include <iostream>
struct instruction {
    char* Name;
    char** Arguments;
    int IgnoreError;
    char* Output;
    char* Error;
};
typedef int (__stdcall *f_callback)(instruction*);
typedef int (__stdcall *f_register)(f_callback);
typedef int (__stdcall *f_execute)(instruction*);


char* pchar(std::string input){
    return const_cast<char*>(input.c_str());
}

int callback(instruction* i){
    std::cout << "callback of " << i->Name << std::endl;
    std::cout << "output: " << i->Output << std::endl;
    std::cout << "error: " << i->Error << std::endl;
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
    registerHandler(callback);

    char* name = pchar("ins from C++");
    char** args = new char* [2];
    args[0] = pchar("motor");
    args[1] = pchar("10");

    instruction i = {
        .Name = name,
        .Arguments = args,
        .IgnoreError = 0,
    };
    //std::cout << "--testing--" << std::endl;
    //callback(&i);
    //std::cout << "--executing--" << std::endl;
    execute(&i);
}
