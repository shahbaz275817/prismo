package utils

import (
	"testing"
)

func TestCreateAWBLegTypeLockKey(t *testing.T) {
	tests := []struct {
		awbNumber string
		legType   string
		expected  string
	}{
		{
			awbNumber: "AWB123",
			legType:   "first_mile",
			expected:  "AWB123first_mile",
		},
	}
	for _, test := range tests {
		result := CreateAWBLegTypeLockKey(test.awbNumber, test.legType)

		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}
