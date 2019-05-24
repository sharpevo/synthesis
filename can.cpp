#include <iostream>
#include <windows.h>
//#include  "can.h"
//#include <windows.h>
//extern "C" {
//}
typedef int (__stdcall *f_funci)(char*, char**);

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
    return 0;
}
