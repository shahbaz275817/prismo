package utils

import "testing"

func TestFormatFloat32IntoString(t *testing.T) {
	tests := []struct {
		Input     float32
		ExpOutput string
	}{
		{
			Input:     0.0,
			ExpOutput: "0",
		},

		{
			Input:     0.1,
			ExpOutput: "0.1",
		},

		{
			Input:     0.111,
			ExpOutput: "0.111",
		},

		{
			Input:     1.123,
			ExpOutput: "1.123",
		},
	}

	for _, test := range tests {
		t.Run("TestFormatFloat32IntoString", func(t *testing.T) {

			output := FormatFloat32IntoString(test.Input)

			if output != test.ExpOutput {
				t.Errorf("got %s | wanted %s", output, test.ExpOutput)
			}

		})
	}
}
