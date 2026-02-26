//go:build cuda

package quantx

import (
	"testing"
	"unsafe"
)

func BenchmarkDequantizeCUDA(b *testing.B) {
	sizes := []struct {
		name string
		k    int
	}{
		{"1MB", 1024 * 256 / 4},
		{"4MB", 1024 * 1024 / 4},
		{"16MB", 16 * 1024 * 1024 / 4},
		{"64MB", 64 * 1024 * 1024 / 4},
		{"256MB", 256 * 1024 * 1024 / 4},
		{"1GB", 1024 * 1024 * 1024 / 4},
	}

	for _, sz := range sizes {
		b.Run("CUDA_"+sz.name, func(b *testing.B) {
			k := sz.k
			numBlocks := k / 256
			if numBlocks == 0 { numBlocks = 1; k = 256 }
			blocks := make([]BlockQ2K, numBlocks)
			for i := range blocks {
				blocks[i].D = 1.0
				blocks[i].Dmin = 0.0
			}

			inputSize := numBlocks * 72
			outputSize := k * 4

			dVx, err := NewCudaBuffer(inputSize)
			if err != nil { b.Skip("CUDA not available:", err); return }
			defer dVx.Free()

			dVy, err := NewCudaBuffer(outputSize)
			if err != nil { b.Skip("CUDA not available:", err); return }
			defer dVy.Free()

			dVx.CopyToDevice(unsafe.Pointer(&blocks[0]), inputSize)

			b.SetBytes(int64(k * 4))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				DequantizeQ2KCUDAKernel(dVx.Ptr, dVy.Ptr, k)
			}
		})
	}
}
