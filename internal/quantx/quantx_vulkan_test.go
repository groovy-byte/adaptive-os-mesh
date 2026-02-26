//go:build vulkan

package quantx

import (
	"testing"
	"unsafe"
)

func BenchmarkDequantizeVulkan(b *testing.B) {
	sizes := []struct {
		name string
		k    int
	}{
		{"1MB", 1024 * 256 / 4},
		{"4MB", 1024 * 1024 / 4},
		{"16MB", 16 * 1024 * 1024 / 4},
		{"64MB", 64 * 1024 * 1024 / 4},
		{"256MB", 256 * 1024 * 1024 / 4},
	}

	for _, sz := range sizes {
		b.Run("Vulkan_"+sz.name, func(b *testing.B) {
			k := sz.k
			numBlocks := k / 256
			if numBlocks == 0 { numBlocks = 1; k = 256 }
			blocks := make([]BlockQ2K, numBlocks)
			for i := range blocks {
				blocks[i].D = 1.0
				blocks[i].Dmin = 0.0
			}
			output := make([]float32, k)
			
			err := DequantizeQ2KVulkanPrepare(k)
			if err != nil { b.Skip("Vulkan not available:", err); return }

			ptr := unsafe.Pointer(&blocks[0])

			b.SetBytes(int64(k * 4))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				DequantizeQ2KVulkan(ptr, output, k)
			}
		})
	}
}

func BenchmarkDequantizeVulkanKernel(b *testing.B) {
	sizes := []struct {
		name string
		k    int
	}{
		{"1MB", 1024 * 256 / 4},
		{"4MB", 1024 * 1024 / 4},
		{"16MB", 16 * 1024 * 1024 / 4},
		{"64MB", 64 * 1024 * 1024 / 4},
		{"256MB", 256 * 1024 * 1024 / 4},
	}

	for _, sz := range sizes {
		b.Run("VulkanKernel_"+sz.name, func(b *testing.B) {
			k := sz.k
			if k < 256 { k = 256 }
			
			err := DequantizeQ2KVulkanPrepare(k)
			if err != nil { b.Skip("Vulkan not available:", err); return }

			b.SetBytes(int64(k * 4))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				DequantizeQ2KVulkanKernel(k)
			}
		})
	}
}
