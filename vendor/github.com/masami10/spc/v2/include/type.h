//
// Created by 华为 on 2022/3/14.
//

#ifndef SPC_TYPE_H
#define SPC_TYPE_H

#include "common.h"

#ifdef __cplusplus
extern "C" {
#endif


typedef struct {
    int ret;
    double data;
} CALC_RET;

typedef struct {
    double *pData;
    size_t lData;
    double center;
    double lower;
    double upper;
} SPC_RET;

typedef struct {
    double *pXData;
    size_t lXData;
    double *pYData;
    size_t lYData;
} PLOT_RET;


int memcpy_plot_ret(PLOT_RET **ret, const double *pX, size_t lX, const double *pY, size_t lY);

void free_plot_ret(PLOT_RET **ret);

int memcpy_spc_ret(SPC_RET **ret, const double *p, size_t size);

void free_spc_ret(SPC_RET **ret);

#ifdef __cplusplus
}
#endif

#endif //SPC_TYPE_H
