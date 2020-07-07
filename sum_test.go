package gotesting

import "testing"

func TestSum(t *testing.T) {
	sum := Sum(1, 1)

	if sum != 2 {
		t.Errorf("Sum(1, 1) = %d; want 2", sum)
	}
}
