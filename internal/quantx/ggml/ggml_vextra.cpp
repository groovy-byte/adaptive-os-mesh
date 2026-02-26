#include "ggml_vextra.h"
#include "ggml-backend-impl.h"
#include "ggml-impl.h"
#include "../quantx.h"
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <vector>

// --- ScheInfer C Interface ---
extern "C" char * scheinfer_route_task(uint64_t data_size_bytes);

// --- Vextra Buffer Type ---

static const char * ggml_backend_vextra_buffer_type_get_name(ggml_backend_buffer_type_t buft) {
    (void)buft;
    return "Vextra";
}

static struct ggml_backend_buffer_i ggml_backend_vextra_buffer_interface = {
    /* .free_buffer  = */ [](ggml_backend_buffer_t b) { free(b->context); },
    /* .get_base     = */ [](ggml_backend_buffer_t b) { return b->context; },
    /* .init_tensor  = */ nullptr,
    /* .memset_tensor = */ nullptr,
    /* .set_tensor   = */ nullptr,
    /* .get_tensor   = */ nullptr,
    /* .cpy_tensor   = */ nullptr,
    /* .clear        = */ nullptr,
    /* .reset        = */ nullptr,
};

static ggml_backend_buffer_t ggml_backend_vextra_buffer_type_alloc_buffer(ggml_backend_buffer_type_t buft, size_t size) {
    void * data = malloc(size);
    if (!data) return nullptr;
    return ggml_backend_buffer_init(buft, ggml_backend_vextra_buffer_interface, data, size);
}

static struct ggml_backend_buffer_type_i ggml_backend_vextra_buffer_type_interface = {
    /* .get_name         = */ ggml_backend_vextra_buffer_type_get_name,
    /* .alloc_buffer     = */ ggml_backend_vextra_buffer_type_alloc_buffer,
    /* .get_alignment    = */ [](ggml_backend_buffer_type_t) { return (size_t)32; },
    /* .get_max_size     = */ nullptr,
    /* .get_alloc_size   = */ nullptr,
    /* .is_host          = */ [](ggml_backend_buffer_type_t) { return true; },
};

static struct ggml_backend_buffer_type ggml_backend_vextra_buffer_type_struct = {
    /* .iface   = */ ggml_backend_vextra_buffer_type_interface,
    /* .device  = */ nullptr,
    /* .context = */ nullptr,
};

ggml_backend_buffer_type_t ggml_backend_vextra_buffer_type() {
    return &ggml_backend_vextra_buffer_type_struct;
}

// --- Vextra Backend ---

struct ggml_backend_vextra_context {
    int placeholder;
};

static const char * ggml_backend_vextra_get_name(ggml_backend_t backend) {
    (void)backend;
    return "Vextra";
}

static void ggml_backend_vextra_free(ggml_backend_t backend) {
    delete (ggml_backend_vextra_context *)backend->context;
    delete backend;
}

static enum ggml_status ggml_backend_vextra_graph_compute(ggml_backend_t backend, struct ggml_cgraph * cgraph) {
    (void)backend;
    for (int i = 0; i < cgraph->n_nodes; i++) {
        struct ggml_tensor * node = cgraph->nodes[i];
        
        // 1. ScheInfer Routing Decision
        if (node->op == GGML_OP_MUL_MAT) {
            size_t data_size = ggml_nelements(node->src[0]) * ggml_type_size(node->src[0]->type);
            char * provider = scheinfer_route_task(data_size);
            
            // For now we just log the decision in this mock
            // printf("[Vextra] node %s size %zu -> %s\n", node->name, data_size, provider);
            
            free(provider);
        }

        // 2. Intercept CPY for dequantization
        if (node->op == GGML_OP_CPY) {
            struct ggml_tensor * src = node->src[0];
            struct ggml_tensor * dst = node; 
            
            if (src->type == GGML_TYPE_Q2_K && dst->type == GGML_TYPE_F32) {
                int k = ggml_nelements(src);
                dequantize_q2_k_avx2(src->data, (float *)dst->data, k);
                continue;
            }
        }
    }
    return GGML_STATUS_SUCCESS;
}

static struct ggml_backend_i ggml_backend_vextra_interface = {
    /* .get_name                = */ ggml_backend_vextra_get_name,
    /* .free                    = */ ggml_backend_vextra_free,
    /* .set_tensor_async        = */ nullptr,
    /* .get_tensor_async        = */ nullptr,
    /* .cpy_tensor_async        = */ nullptr,
    /* .synchronize             = */ [](ggml_backend_t) {},
    /* .graph_plan_create       = */ nullptr,
    /* .graph_plan_free         = */ nullptr,
    /* .graph_plan_update       = */ nullptr,
    /* .graph_plan_compute      = */ nullptr,
    /* .graph_compute           = */ ggml_backend_vextra_graph_compute,
    /* .event_record            = */ nullptr,
    /* .event_wait              = */ nullptr,
    /* .graph_optimize          = */ nullptr,
};

ggml_backend_t ggml_backend_vextra_init() {
    ggml_backend_t backend = new ggml_backend;
    std::memset(backend, 0, sizeof(struct ggml_backend));
    backend->iface = ggml_backend_vextra_interface;
    backend->context = new ggml_backend_vextra_context();
    return backend;
}
