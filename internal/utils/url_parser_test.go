package utils

import (
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestURLParser_Get(t *testing.T) {
	u, _ := url.Parse("http://example.com/?param1=value1&param2=value2")
	up := NewURLParser(u)

	value := up.Get("param1")
	if value != "value1" {
		t.Errorf("Get() failed, expected 'value1', got '%s'", value)
	}

	value = up.Get("nonexistent")
	if value != "" {
		t.Errorf("Get() failed, expected empty string, got '%s'", value)
	}
}

func TestURLParser_GetList(t *testing.T) {
	u, _ := url.Parse("http://example.com/?param1=value1,value2,value3")
	up := NewURLParser(u)

	values := up.GetList("param1")
	expected := []string{"value1", "value2", "value3"}
	if !reflect.DeepEqual(values, expected) {
		t.Errorf("GetList() failed, expected %v, got %v", expected, values)
	}

	values = up.GetList("nonexistent")
	if values != nil {
		t.Errorf("GetList() failed, expected empty slice, got %v", values)
	}
}

func TestURLParser_GetCreatedRangeParams(t *testing.T) {
	testCases := []struct {
		name                string
		url                 string
		expectedCreatedFrom *time.Time
		expectedCreatedTo   *time.Time
		expectedErr         bool
		expectedErrText     string
	}{
		{
			name:                "Valid parameters",
			url:                 "http://example.com/?created_at_from=1619500800&created_at_to=1619587200",
			expectedCreatedFrom: getTimePointer(time.Unix(1619500800, 0).UTC()),
			expectedCreatedTo:   getTimePointer(time.Unix(1619587200, 0).UTC()),
			expectedErr:         false,
			expectedErrText:     "",
		},
		{
			name:                "Missing 'created_at_from' parameter",
			url:                 "http://example.com/?created_at_to=1619587200",
			expectedCreatedFrom: nil,
			expectedCreatedTo:   nil,
			expectedErr:         true,
			expectedErrText:     "created_at_from is missing",
		},
		{
			name:                "Missing 'created_at_to' parameter",
			url:                 "http://example.com/?created_at_from=1619500800",
			expectedCreatedFrom: nil,
			expectedCreatedTo:   nil,
			expectedErr:         true,
			expectedErrText:     "created_at_to is missing",
		},
		{
			name:                "Invalid 'created_at_from' value",
			url:                 "http://example.com/?created_at_from=invalid&created_at_to=1619587200",
			expectedCreatedFrom: nil,
			expectedCreatedTo:   nil,
			expectedErr:         true,
			expectedErrText:     "created_at_from should be unix epoch",
		},
		{
			name:                "Invalid 'created_at_to' value",
			url:                 "http://example.com/?created_at_from=1619500800&created_at_to=invalid",
			expectedCreatedFrom: nil,
			expectedCreatedTo:   nil,
			expectedErr:         true,
			expectedErrText:     "created_at_to should be unix epoch",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse(tc.url)
			up := NewURLParser(u)

			createdFrom, createdTo, err := up.GetCreatedRangeParams()
			if (err != nil) != tc.expectedErr {
				t.Errorf("GetCreatedRangeParams() error mismatch for test case '%s'. Expected error: %v, Got error: %v", tc.name, tc.expectedErr, err != nil)
			}
			if tc.expectedErr && err.Error() != tc.expectedErrText {
				t.Errorf("GetCreatedRangeParams() error text mismatch for test case '%s'. Expected error: '%s', Got error: '%s'", tc.name, tc.expectedErrText, err.Error())
			}
			if !timesEqual(createdFrom, tc.expectedCreatedFrom) || !timesEqual(createdTo, tc.expectedCreatedTo) {
				t.Errorf("GetCreatedRangeParams() failed for test case '%s'. Unexpected values. Expected: createdFrom=%v, createdTo=%v, Got: createdFrom=%v, createdTo=%v", tc.name, tc.expectedCreatedFrom, tc.expectedCreatedTo, createdFrom, createdTo)
			}
		})
	}
}

func TestURLParser_GetInboundedRangeParams(t *testing.T) {
	testCases := []struct {
		name                  string
		url                   string
		expectedInboundedFrom *time.Time
		expectedInboundedTo   *time.Time
		expectedErr           bool
		expectedErrText       string
	}{
		{
			name:                  "Valid parameters",
			url:                   "http://example.com/?inbounded_at_from=1619500800&inbounded_at_to=1619587200",
			expectedInboundedFrom: getTimePointer(time.Unix(1619500800, 0).UTC()),
			expectedInboundedTo:   getTimePointer(time.Unix(1619587200, 0).UTC()),
			expectedErr:           false,
			expectedErrText:       "",
		},
		{
			name:                  "Missing 'inbounded_at_from' parameter",
			url:                   "http://example.com/?inbounded_at_to=1619587200",
			expectedInboundedFrom: nil,
			expectedInboundedTo:   nil,
			expectedErr:           true,
			expectedErrText:       "inbounded_at_from is missing",
		},
		{
			name:                  "Missing 'inbounded_at_to' parameter",
			url:                   "http://example.com/?inbounded_at_from=1619500800",
			expectedInboundedFrom: nil,
			expectedInboundedTo:   nil,
			expectedErr:           true,
			expectedErrText:       "inbounded_at_to is missing",
		},
		{
			name:                  "Invalid 'inbounded_at_from' value",
			url:                   "http://example.com/?inbounded_at_from=invalid&inbounded_at_to=1619587200",
			expectedInboundedFrom: nil,
			expectedInboundedTo:   nil,
			expectedErr:           true,
			expectedErrText:       "inbounded_at_from should be unix epoch",
		},
		{
			name:                  "Invalid 'inbounded_at_to' value",
			url:                   "http://example.com/?inbounded_at_from=1619500800&inbounded_at_to=invalid",
			expectedInboundedFrom: nil,
			expectedInboundedTo:   nil,
			expectedErr:           true,
			expectedErrText:       "inbounded_at_to should be unix epoch",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse(tc.url)
			up := NewURLParser(u)

			inboundedFrom, inboundedTo, err := up.GetInboundedRangeParams()
			if (err != nil) != tc.expectedErr {
				t.Errorf("GetInboundedRangeParams() error mismatch for test case '%s'. Expected error: %v, Got error: %v", tc.name, tc.expectedErr, err != nil)
			}
			if tc.expectedErr && err.Error() != tc.expectedErrText {
				t.Errorf("GetInboundedRangeParams() error text mismatch for test case '%s'. Expected error: '%s', Got error: '%s'", tc.name, tc.expectedErrText, err.Error())
			}
			if !timesEqual(inboundedFrom, tc.expectedInboundedFrom) || !timesEqual(inboundedTo, tc.expectedInboundedTo) {
				t.Errorf("GetInboundedRangeParams() failed for test case '%s'. Unexpected values. Expected: createdFrom=%v, createdTo=%v, Got: createdFrom=%v, createdTo=%v", tc.name, tc.expectedInboundedFrom, tc.expectedInboundedTo, inboundedFrom, inboundedTo)
			}
		})
	}
}

func TestURLParser_GetRpuAtRangeParams(t *testing.T) {
	testCases := []struct {
		name            string
		url             string
		expectedRpuFrom *time.Time
		expectedRpuTo   *time.Time
		expectedErr     bool
		expectedErrText string
	}{
		{
			name:            "Valid parameters",
			url:             "http://example.com/?rpu_at_from=1619500800&rpu_at_to=1619587200",
			expectedRpuFrom: getTimePointer(time.Unix(1619500800, 0).UTC()),
			expectedRpuTo:   getTimePointer(time.Unix(1619587200, 0).UTC()),
			expectedErr:     false,
			expectedErrText: "",
		},
		{
			name:            "Missing 'rpu_at_from' parameter",
			url:             "http://example.com/?rpu_at_to=1619587200",
			expectedRpuFrom: nil,
			expectedRpuTo:   nil,
			expectedErr:     true,
			expectedErrText: "rpu_at_from is missing",
		},
		{
			name:            "Missing 'rpu_at_to' parameter",
			url:             "http://example.com/?rpu_at_from=1619500800",
			expectedRpuFrom: nil,
			expectedRpuTo:   nil,
			expectedErr:     true,
			expectedErrText: "rpu_at_to is missing",
		},
		{
			name:            "Invalid 'rpu_at_from' value",
			url:             "http://example.com/?rpu_at_from=invalid&rpu_at_to=1619587200",
			expectedRpuFrom: nil,
			expectedRpuTo:   nil,
			expectedErr:     true,
			expectedErrText: "rpu_at_from should be unix epoch",
		},
		{
			name:            "Invalid 'rpu_at_to' value",
			url:             "http://example.com/?rpu_at_from=1619500800&rpu_at_to=invalid",
			expectedRpuFrom: nil,
			expectedRpuTo:   nil,
			expectedErr:     true,
			expectedErrText: "rpu_at_to should be unix epoch",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse(tc.url)
			up := NewURLParser(u)

			rpuFrom, rpuTo, err := up.GetRpuAtRangeParams()
			if (err != nil) != tc.expectedErr {
				t.Errorf("GetCreatedRangeParams() error mismatch for test case '%s'. Expected error: %v, Got error: %v", tc.name, tc.expectedErr, err != nil)
			}
			if tc.expectedErr && err.Error() != tc.expectedErrText {
				t.Errorf("GetCreatedRangeParams() error text mismatch for test case '%s'. Expected error: '%s', Got error: '%s'", tc.name, tc.expectedErrText, err.Error())
			}
			if !timesEqual(rpuFrom, tc.expectedRpuFrom) || !timesEqual(rpuTo, tc.expectedRpuTo) {
				t.Errorf("GetCreatedRangeParams() failed for test case '%s'. Unexpected values. Expected: createdFrom=%v, createdTo=%v, Got: createdFrom=%v, createdTo=%v", tc.name, tc.expectedRpuFrom, tc.expectedRpuTo, rpuFrom, rpuTo)
			}
		})
	}
}

// Helper function to compare time.Time values while handling nil pointers
func timesEqual(t1, t2 *time.Time) bool {
	if t1 == nil && t2 == nil {
		return true
	}
	if t1 == nil || t2 == nil {
		return false
	}
	return t1.Equal(*t2)
}

func TestURLParser_GetBool(t *testing.T) {
	testCases := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "'true'",
			url:      "http://example.com/?flag=true",
			expected: true,
		},
		{
			name:     "'false'",
			url:      "http://example.com/?flag=false",
			expected: false,
		},
		{
			name:     "Incorrect boolean value key",
			url:      "http://example.com/?flag=incorrectBool",
			expected: false,
		},
		{
			name:     "Nonexistent key",
			url:      "http://example.com/",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse(tc.url)
			up := NewURLParser(u)

			result := up.GetBool("flag")
			if result != tc.expected {
				t.Errorf("GetBool() failed for key '%s', expected %v, got %v", "flag", tc.expected, result)
			}
		})
	}
}

func TestURLParser_GetPaginationParams(t *testing.T) {
	testCases := []struct {
		name            string
		url             string
		expectedLimit   int
		expectedOffset  int
		expectedErr     bool
		expectedErrText string
	}{
		{
			name:            "Valid parameters",
			url:             "http://example.com/?page=2&per_page=25",
			expectedLimit:   25,
			expectedOffset:  25,
			expectedErr:     false,
			expectedErrText: "",
		},
		{
			name:            "Missing 'page' parameter",
			url:             "http://example.com/?per_page=25",
			expectedLimit:   25,
			expectedOffset:  0,
			expectedErr:     false,
			expectedErrText: "",
		},
		{
			name:            "Missing 'per_page' parameter",
			url:             "http://example.com/?page=2",
			expectedLimit:   defaultLimit,
			expectedOffset:  1 * defaultLimit, //todo fix this test case
			expectedErr:     false,
			expectedErrText: "",
		},
		{
			name:            "Invalid 'page' value",
			url:             "http://example.com/?page=invalid&per_page=25",
			expectedLimit:   0,
			expectedOffset:  0,
			expectedErr:     true,
			expectedErrText: "page should be a number",
		},
		{
			name:            "Invalid 'per_page' value",
			url:             "http://example.com/?page=2&per_page=invalid",
			expectedLimit:   0,
			expectedOffset:  0,
			expectedErr:     true,
			expectedErrText: "per_page should be a number",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse(tc.url)
			up := NewURLParser(u)

			limit, offset, err := up.GetPaginationParams()
			if (err != nil) != tc.expectedErr {
				t.Errorf("GetPaginationParams() error mismatch for test case '%s'. Expected error: %v, Got error: %v", tc.name, tc.expectedErr, err != nil)
			}
			if tc.expectedErr && err.Error() != tc.expectedErrText {
				t.Errorf("GetPaginationParams() error text mismatch for test case '%s'. Expected error: '%s', Got error: '%s'", tc.name, tc.expectedErrText, err.Error())
			}
			if limit != tc.expectedLimit || offset != tc.expectedOffset {
				t.Errorf("GetPaginationParams() failed for test case '%s'. Unexpected values. Expected: limit=%d, offset=%d, Got: limit=%d, offset=%d", tc.name, tc.expectedLimit, tc.expectedOffset, limit, offset)
			}
		})
	}
}

func getTimePointer(t time.Time) *time.Time {
	return &t
}
