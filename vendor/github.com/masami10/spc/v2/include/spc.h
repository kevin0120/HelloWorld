//
// Created by 华为 on 2022/3/11.
//

#ifndef SPC_SPC_H
#define SPC_SPC_H

#include "common.h"
#include "error.h"

#ifdef __cplusplus
extern "C" {
#endif

double cpk(double *data, size_t length, double usl, double lsl);

double cmk(double *data, size_t length, double usl, double lsl);

double cp(double *data, size_t length, double usl, double lsl);

double cr(double *data, size_t length, double usl, double lsl);

double calc_std(double *data, size_t length);

#ifdef __cplusplus
}
#endif

ST_RET xbarSbar(double **data, size_t n_col, size_t n_row, SPC_RET **spc_ret);


#endif //SPC_SPC_H
