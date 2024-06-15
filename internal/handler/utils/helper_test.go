package utils

import (
	"testing"
	"time"
)

func TestFormatVolumeFloat32(t *testing.T) {
	tests := []struct {
		input    float32
		expected string
	}{

		{
			input:    2.12,
			expected: "2.120000 m^3",
		},

		{
			input:    2.333000,
			expected: "2.333000 m^3",
		},

		// %f default precision is upto 6 places after decimal
		{
			input:    2.12345678,
			expected: "2.123457 m^3",
		},
	}

	for _, test := range tests {
		t.Run("TestFormatVolumeFloat32", func(t *testing.T) {
			output := FormatVolumeFloat32(&test.input)

			if output != test.expected {
				t.Errorf("Got %s, Wanted %s", output, test.expected)
			}
		})
	}

}

func TestGetUnixMilliTimestampIfPresent(t *testing.T) {
	cases := []struct {
		name          string
		inputTime     *time.Time
		expectedValue *int64
	}{
		{
			name:          "NilInput",
			inputTime:     nil,
			expectedValue: nil,
		},
		{
			name: "ValidInput",
			inputTime: func() *time.Time {
				now := time.Now()
				return &now
			}(),
			expectedValue: func() *int64 {
				now := time.Now()
				expected := now.UnixMilli()
				return &expected
			}(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetUnixMilliTimestampIfPresent(tc.inputTime)
			if result == nil && tc.expectedValue != nil {
				t.Errorf("Expected non-nil result, got nil")
			} else if result != nil && tc.expectedValue == nil {
				t.Errorf("Expected nil result, got %v", *result)
			} else if result != nil && *result != *tc.expectedValue {
				t.Errorf("Expected %d, got %d", *tc.expectedValue, *result)
			}
		})
	}
}

func TestGetIntPointer(t *testing.T) {
	cases := []struct {
		name          string
		inputValue    int64
		expectedValue *int64
	}{
		{
			name:          "PositiveValue",
			inputValue:    123,
			expectedValue: int64Pointer(123),
		},
		{
			name:          "ZeroValue",
			inputValue:    0,
			expectedValue: int64Pointer(0),
		},
		{
			name:          "NegativeValue",
			inputValue:    -456,
			expectedValue: int64Pointer(-456),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetIntPointer(tc.inputValue)
			if result == nil || *result != *tc.expectedValue {
				t.Errorf("Expected %d, got %v", *tc.expectedValue, result)
			}
		})
	}
}

func int64Pointer(value int64) *int64 {
	return &value
}
