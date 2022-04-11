#include <stdio.h>
#include "library.h"

void hello() {
    printf("Hello, World!\n");
}

void PrintArray(double *data, int length) {
    int i;
    double a;
    for (i = 0; i < length; ++i) {
        a = data[i];
        printf("%f, ", data[i]);
    }
    printf("\n");
}