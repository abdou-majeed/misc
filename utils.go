package utils

import (
	"errors"
	"math"
)

// ErrOverflow is returned when an arithmetic operation exceeds integer limits.
var ErrOverflow = errors.New("integer overflow")

// Add safely adds two integers and returns an error if an overflow occurs.
// This is critical for preventing infinite loops or memory panics when 
// checking multiples of very large prime numbers.
func Add(a, b int) (int, error) {
	// Check for positive overflow
	if b > 0 && a > math.MaxInt-b {
		return 0, ErrOverflow
	}
	
	// Check for negative overflow (though unlikely in your prime sieve, it's good practice)
	if b < 0 && a < math.MinInt-b {
		return 0, ErrOverflow
	}
	
	return a + b, nil
}