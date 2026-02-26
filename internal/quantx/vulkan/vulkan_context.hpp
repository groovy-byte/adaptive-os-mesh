//go:build vulkan

#ifndef VULKAN_CONTEXT_HPP
#define VULKAN_CONTEXT_HPP

#ifdef USE_VULKAN
#include <vulkan/vulkan.h>
#include <vector>
#include <stdint.h>

class VulkanContext {
public:
    VulkanContext();
    ~VulkanContext();

    bool init();
    bool prepare(int max_k);
    bool run_kernel(int k);
    bool dequantize(const void* vx, float* vy, int k);

private:
    VkInstance instance;
    VkPhysicalDevice physicalDevice;
    VkDevice device;
    VkQueue computeQueue;
    uint32_t computeQueueFamilyIndex;
    
    VkShaderModule shaderModule;
    VkDescriptorSetLayout descriptorSetLayout;
    VkPipelineLayout pipelineLayout;
    VkPipeline pipeline;
    VkDescriptorPool descriptorSetPool;
    VkDescriptorSet descriptorSet;
    VkCommandPool commandPool;

    // Persistent Buffers
    VkBuffer inputBuffer;
    VkDeviceMemory inputBufferMemory;
    VkBuffer outputBuffer;
    VkDeviceMemory outputBufferMemory;
    int current_max_k;

    bool createInstance();
    bool pickPhysicalDevice();
    bool createDevice();
    bool createDescriptorSetLayout();
    bool createComputePipeline();
    bool createDescriptorSet();
    bool createCommandPool();
    
    uint32_t findMemoryType(uint32_t typeFilter, VkMemoryPropertyFlags properties);
    void destroyBuffers();
};
#endif

#endif
