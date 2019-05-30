#include <windows.h>
#include <iostream>
typedef int (__stdcall *f_callback)(char*, char*);
struct instruction {
    char* Name;
    char** Arguments;
    int IgnoreError;
    f_callback callback;
};
typedef int (__stdcall *f_execute)(instruction*);


char* pchar(std::string input){
    return const_cast<char*>(input.c_str());
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
    char* name = pchar("ins from C++");
    char** args = new char* [2];
    args[0] = pchar("motor");
    args[1] = pchar("10");
    //auto f = [msg](char* a, char* b)->int{
    f_callback f = [](char* a, char* b)->int{
        char* msg = pchar("some updates");
        std::cout << msg << std::endl;
        std::cout << a << "||" << b << std::endl;
        return 22;
    };

    instruction i = {
        name,args,0,f
        //int test(char* a, char* b) {
            //std::cout << msg << std::endl;
            //std::cout << a + " || " + b << std::endl;
            //return 22;
        //},
    };
    std::cout << "--testing--" << std::endl;
    i.callback(args[0], args[1]);
    std::cout << "--executing--" << std::endl;
    execute(&i);
}
