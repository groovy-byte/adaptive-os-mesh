#ifndef GGML_VEXTRA_H
#define GGML_VEXTRA_H

#include "ggml.h"
#include "ggml-backend.h"

#ifdef  __cplusplus
extern "C" {
#endif

    // Initialize the Vextra backend
    // This backend uses QuantX kernels (AVX2/CUDA/Vulkan) for optimized dequantization and matmul
    GGML_API ggml_backend_t ggml_backend_vextra_init();

    // Check if a buffer is a Vextra buffer
    GGML_API bool ggml_backend_is_vextra(ggml_backend_t backend);

    // Get the default buffer type for Vextra
    GGML_API ggml_backend_buffer_type_t ggml_backend_vextra_buffer_type();

#ifdef  __cplusplus
}
#endif

#endif // GGML_VEXTRA_H
