#include <windows.h>
#include <iostream>
typedef int (__stdcall *f_init)(char*, char*, char*, char*, char*, char*,
    char*, char*, char*, char*);
typedef int (__stdcall *f_controlswitcher)(int);
typedef int (__stdcall *f_readhumiture)(double*, double*);

char* pchar(std::string input){
    return const_cast<char*>(input.c_str());
}

int main(){
    HINSTANCE canlib = LoadLibrary("can.dll"); 
    if (!canlib) {
        std::cout << "failed to load dynamic library" << std::endl;
        return 1;
    }
    f_init Init = (f_init)GetProcAddress(canlib, "NewDao");
    if (!Init) {
        std::cout << "failed to load function NewDao" << std::endl;
    }
    f_controlswitcher ControlSwitcher = (f_controlswitcher)GetProcAddress(canlib,
            "ControlSwitcher");
    if (!ControlSwitcher) {
        std::cout << "failed to load function ControlSwitcher" << std::endl;
    }
    f_readhumiture ReadHumiture = (f_readhumiture)GetProcAddress(canlib,
            "ReadHumiture");
    if (!ReadHumiture) {
        std::cout << "failed to load function ReadHumiture" << std::endl;
    }
    if (!Init(pchar("4"),pchar("0"),pchar("0x00000001"),
        pchar("0"), pchar("0x00000000"), pchar("0xFFFFFFFF"), pchar("0"),
        pchar("0x00"), pchar("0x1c"), pchar("0"))){
        std::cout << "failed to init can device" << std::endl;
        return 1;
    }
    if (!ControlSwitcher(12)) {
        std::cout << "failed to send switch instruction" << std::endl;
        return 1;
    }
    double humi, temp;
    if (!ReadHumiture(&humi, &temp)){
        std::cout << "failed to read humiture" << std::endl;
        return 1;
    }
    std::cout << "read humiture from C++: " << humi << std::endl;
    std::cout << "read temperature from C++: " << temp << std::endl;
    return 0;
}

