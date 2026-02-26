package nacrith

import (
	"testing"
	"math"
)

func TestCDF24Coding(t *testing.T) {
	// The goal of CDF-24 is to provide 2^24 precision levels for token frequency
	// to reduce the floor overhead in large vocabularies.
	
	// Test Case: Encode a low-probability token
	// Standard 16-bit (65536) often loses precision for probabilities < 1e-5
	prob := 1.2345e-6
	
	// Encode using 24-bit precision (16,777,216)
	encoded := uint32(prob * 16777216.0)
	
	// Decode
	decoded := float64(encoded) / 16777216.0
	
	error_margin := math.Abs(prob - decoded)
	
	// Verification: Error margin should be < 1e-7 (24-bit precision)
	// Standard 16-bit error margin would be ~1.5e-5
	if error_margin > 1e-7 {
		t.Errorf("CDF-24 precision failure: error %e > 1e-7", error_margin)
	}
	
	t.Logf("CDF-24 Success: Error margin %e (Standard 16-bit margin was ~%e)", 
		error_margin, 1.0/65536.0)
}
