#ifndef QUANTX_H
#define QUANTX_H

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

/**
 * dequantize_q2_k_avx2 dequantizes Q2_K blocks to float32 using AVX2 and FMA.
 * @param vx Pointer to the quantized 2-bit blocks.
 * @param vy Pointer to the output float32 buffer.
 * @param k Number of elements to dequantize.
 */
void dequantize_q2_k_avx2(const void* vx, float* vy, int k);

/**
 * dequantize_q2_k_avx512 dequantizes Q2_K blocks to float32 using AVX512.
 * @param vx Pointer to the quantized 2-bit blocks.
 * @param vy Pointer to the output float32 buffer.
 * @param k Number of elements to dequantize.
 */
void dequantize_q2_k_avx512(const void* vx, float* vy, int k);

/**
 * dequantize_q2_k_cuda dequantizes Q2_K blocks to float32 using CUDA.
 * @param vx Pointer to the quantized 2-bit blocks (on device).
 * @param vy Pointer to the output float32 buffer (on device).
 * @param k Number of elements to dequantize.
 */
void dequantize_q2_k_cuda(const void* vx, float* vy, int k);

#ifdef __cplusplus
}
#endif

#endif // QUANTX_H
