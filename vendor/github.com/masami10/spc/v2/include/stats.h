//
// Created by 华为 on 2022/3/14.
//

#ifndef SPC_STATS_H
#define SPC_STATS_H

#include "common.h"

#ifdef __cplusplus
extern "C" {
#endif

ST_RET CalcStandardDeviation(double *data, size_t length, int ddof, double mean, CALC_RET *ret);

ST_RET CalcMean(double *data, size_t length, CALC_RET *ret);

ST_RET CalcAvg(double *data, size_t length, CALC_RET *ret);

ST_RET CalcHistogram(double *data, size_t length, float usl, float lsl, int step, bool density, PLOT_RET **ret);

ST_RET CalcNormalDist(double *data, size_t length, float usl, float lsl, int step, bool density, PLOT_RET **ret);

double normFun(double x, double mu, double sigma);
void PrintArray(double *data,int length);
#ifdef __cplusplus
}
#endif

#endif //SPC_STATS_H
