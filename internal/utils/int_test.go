package utils

import (
	"reflect"
	"testing"
)

func TestConvertSliceToInt(t *testing.T) {
	tests := []struct {
		input       []string
		expected    []int
		errExpected bool
	}{
		{
			input:       []string{"1", "2", "3", "4", "5"},
			expected:    []int{1, 2, 3, 4, 5},
			errExpected: false,
		},
		{
			input:       []string{"1", "2", "3", "4", "5", "abc"},
			expected:    nil,
			errExpected: true,
		},
		{
			input:       []string{},
			expected:    []int{},
			errExpected: false,
		},
	}
	for _, test := range tests {
		result, err := ConvertSliceOfStringToSliceOfInt(test.input)
		if test.errExpected {
			if err == nil {
				t.Errorf("Error expected but not received for input %v", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input %v: %v", test.input, err)
			}
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Unexpected result for input %v: %v (expected %v)", test.input, result, test.expected)
			}
		}
	}
}

func TestConvertSliceToInt32(t *testing.T) {
	tests := []struct {
		input       []string
		expected    []int32
		errExpected bool
	}{
		{
			input:       []string{"1", "2", "3", "4", "5"},
			expected:    []int32{1, 2, 3, 4, 5},
			errExpected: false,
		},
		{
			input:       []string{"1", "2", "3", "4", "5", "abc"},
			expected:    nil,
			errExpected: true,
		},
		{
			input:       []string{},
			expected:    []int32{},
			errExpected: false,
		},
	}
	for _, test := range tests {
		result, err := ConvertSliceOfStringToSliceOfInt32(test.input)
		if test.errExpected {
			if err == nil {
				t.Errorf("Error expected but not received for input %v", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input %v: %v", test.input, err)
			}
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Unexpected result for input %v: %v (expected %v)", test.input, result, test.expected)
			}
		}
	}
}
