#!/bin/bash
x86_64-w64-mingw32-g++ -static-libgcc -static-libstdc++ -Wl,-Bstatic -lstdc++ -lpthread -Wl,-Bdynamic can.cpp
