#!/bin/bash
x86_64-w64-mingw32-g++ -static-libgcc -static-libstdc++ -Wl,-Bstatic,--whole-archive -lwinpthread -Wl,--no-whole-archive -o interpreter.exe interpreter.cpp
