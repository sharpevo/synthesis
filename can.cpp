#include <iostream>
#include <windows.h>
typedef int (__stdcall *f_init)(char*, char*, char*, char*, char*, char*,
    char*, char*, char*, char*);
typedef int (__stdcall *f_readhumiture)(double*, double*);

char* conv(std::string input){
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
    f_readhumiture ReadHumiture = (f_readhumiture)GetProcAddress(canlib, "ReadHumiture");
    if (!ReadHumiture) {
        std::cout << "could not locate the function ReadHumiture" << std::endl;
    }

    std::string devtype, devindex, devid, canindex, acccode, accmask, filter,
        timing0, timing1, mode;

    devtype = "4";
    devindex = "0";
    devid = "0x00000001";
    canindex = "0";
    acccode = "0x00000000";
    accmask = "0xFFFFFFFF";
    filter = "0";
    timing0 = "0x00";
    timing1 = "0x1c";
    mode = "0";

    std::cout << "Init: " << Init(conv(devtype),conv(devindex),conv(devid),
        conv(canindex), conv(acccode), conv(accmask), conv(filter),
        conv(timing0), conv(timing1), conv(mode)) << std::endl;

    double humi, temp;
    std::cout << "ReadHumiture: " << ReadHumiture(&humi, &temp) << std::endl;
    std::cout << "humiture: " << humi  << std::endl;
    std::cout << "temperature: " << temp  << std::endl;

    return 0;
}

