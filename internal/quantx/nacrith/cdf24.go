package nacrith

import (
	"fmt"
)

// CDF24 handles high-precision arithmetic coding for vocabulary distribution
type CDF24 struct {
	precision uint32
}

func NewCDF24() *CDF24 {
	return &CDF24{
		precision: 1 << 24, // 16,777,216 levels
	}
}

// Encode converts a probability to a 24-bit fixed point representation
func (c *CDF24) Encode(prob float64) (uint32, error) {
	if prob < 0 || prob > 1.0 {
		return 0, fmt.Errorf("invalid probability: %f", prob)
	}
	return uint32(prob * float64(c.precision)), nil
}

// Decode converts a 24-bit fixed point value back to a float64 probability
func (c *CDF24) Decode(value uint32) float64 {
	return float64(value) / float64(c.precision)
}

// BatchEncode scales a whole frequency vector to 24-bit
func (c *CDF24) BatchEncode(probs []float64) []uint32 {
	encoded := make([]uint32, len(probs))
	for i, p := range probs {
		encoded[i], _ = c.Encode(p)
	}
	return encoded
}
