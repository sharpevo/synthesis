#include <iostream>
#include <windows.h>
//#include  "can.h"
//#include <windows.h>
//extern "C" {
//}
typedef int (__stdcall *f_funci)(char*, char**);
typedef int (__stdcall *f_readhumiture)(double*, double*);
typedef int (__stdcall *f_newdao)(char*, char*,char*,char*,char*,char*,char*,char*,char*,char*);

char* conv(std::string input){
    return const_cast<char*>(input.c_str());
}

int main(){
    //Test();
    std::string input = "0123456";
    char* c = const_cast<char*>(input.c_str());
    std::cout << c << std::endl;

    //std::string output;
    //char* ch = new char[output.size()+1];
    //std::copy(output.begin(), output.end(), ch);
    //ch[output.size()] = '\0';

    HINSTANCE a = LoadLibrary("can.dll"); 
    if (!a) {
        std::cout << "could not load the dynamic library" << std::endl;
        return 1;
    }
    f_funci Test = (f_funci)GetProcAddress(a, "Test");
    if (!Test) {
        std::cout << "could not locate the function" << std::endl;
    }

    //Test(c, &c);
    std::cout <<"---- start" << std::endl;
    std::cout << "Test() returned: " << Test(c, &c) << std::endl;
    std::cout << "output: " << c << std::endl; // TODO: utf8
    std::cout <<"---- done"<< std::endl;


    f_newdao NewDao = (f_newdao)GetProcAddress(a, "NewDao");
    if (!NewDao) {
        std::cout << "could not locate the function: "<< "NewDao" << std::endl;
    }

    f_readhumiture ReadHumiture = (f_readhumiture)GetProcAddress(a, "ReadHumiture");
    if (!ReadHumiture) {
        std::cout << "could not locate the function: "<< "ReadHumiture" << std::endl;
    }

    std::string devtype = "4";
    std::string devindex = "0";
    std::string devid = "0x00000001";
    std::string canindex = "0";
    std::string acccode = "0x00000000";
    std::string accmask = "0xFFFFFFFF";
    std::string filter = "0";
    std::string timing0 = "0x00";
    std::string timing1 = "0x1c";
    std::string mode = "0";

    int dao = NewDao(conv(devtype),conv(devindex),conv(devid), conv(canindex), conv(acccode), conv(accmask), conv(filter), conv(timing0), conv(timing1), conv(mode));
    std::cout << "NewDao: " << dao  << std::endl;


    //double* humi = 0;
    //double* temp = 0;
    double humi;
    double temp;
    int humiture = ReadHumiture(&humi, &temp);
    std::cout << "ReadHumiture: " << humiture  << std::endl;
    std::cout << "humiture: " << humi  << std::endl;
    std::cout << "temperature: " << temp  << std::endl;

    return 0;
}

