package main

/*
#cgo CFLAGS: -I ${SRCDIR}/../include/
#include "library.h"

#include <stdio.h>
#include <windows.h>


typedef void (CALLBACK *AddFunc)();
void hello2() {
      HMODULE hDll = LoadLibrary("libspc_lib.dll");

      //printf("Library initialized failed:%d\n",1);
      if (hDll != NULL)
      {
           AddFunc add = (AddFunc)GetProcAddress(hDll, "hello");
           //printf("Library initialized failed:%d\n",2);
            if (add)
            {
                  add();
                  //printf("Library initialized failed:%d\n",3);
            }
            FreeLibrary(hDll);
			 //printf("Library initialized failed:%d\n",4);
      }
}

*/
import "C"

import (
	_ "HelloWorld/test_CGO/CHello/MainDll/driver"
)

func main() {
	//fmt.Println("hh")
	C.hello2()
	select {}
}
