package quantx

import (
	"testing"
	"unsafe"
	"math"
)

type blockQ2K struct {
	D    float32
	Dmin float32
	Qs   [64]uint8
}

func TestDequantizeQ2K(t *testing.T) {
	// Prepare a single block (256 elements)
	block := blockQ2K{
		D:    2.0,
		Dmin: 0.5,
	}

	// Pattern: 0, 1, 2, 3 repeating
	// Each byte holds 4 elements: (3 << 6) | (2 << 4) | (1 << 2) | (0 << 0) = 0xE4
	// In our code: v_shifts = [6, 4, 2, 0]. So bit 0-1 is element 0, bit 2-3 is element 1...
	pattern := uint8(0xE4)
	for i := range block.Qs {
		block.Qs[i] = pattern
	}

	k := 256
	output := make([]float32, k)
	
	DequantizeQ2K(unsafe.Pointer(&block), output, k)

	expected := []float32{0.5, 2.5, 4.5, 6.5} // (0*2+0.5, 1*2+0.5, 2*2+0.5, 3*2+0.5)
	
	for i := 0; i < k; i++ {
		val := output[i]
		exp := expected[i%4]
		if math.Abs(float64(val-exp)) > 1e-5 {
			t.Errorf("At index %d: expected %f, got %f", i, exp, val)
		}
	}
}

func BenchmarkDequantize(b *testing.B) {
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
		b.Run("CPU_"+sz.name, func(b *testing.B) {
			k := sz.k
			numBlocks := k / 256
			if numBlocks == 0 { numBlocks = 1; k = 256 }
			blocks := make([]BlockQ2K, numBlocks)
			for i := range blocks {
				blocks[i].D = 1.0
				blocks[i].Dmin = 0.0
			}
			output := make([]float32, k)
			ptr := unsafe.Pointer(&blocks[0])

			b.SetBytes(int64(k * 4))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				DequantizeQ2K(ptr, output, k)
			}
		})
	}
}
