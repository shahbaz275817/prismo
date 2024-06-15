package utils

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextLogger(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/endpoint", nil)

	ctx, lgr := ContextLogger(req)

	assert.NotNil(t, ctx)
	lgr.Infof("Some logs - check console for verification")
}

func TestContextLoggerURLParser(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/endpoint?name=amphibian", nil)

	ctx, lgr, urlParser := ContextLoggerURLParser(req)

	assert.NotNil(t, ctx)
	lgr.Infof("Some logs - check console for verification")
	assert.NotNil(t, urlParser)
	assert.Equal(t, "amphibian", urlParser.Get("name"))
}

func TestGetRouteIndexFromParam(t *testing.T) {
	testCases := []struct {
		name           string
		routeSeqArray  []string
		expectedResult []int
		expectedError  error
	}{
		{
			name:           "Valid route index array",
			routeSeqArray:  []string{"1", "2", "3"},
			expectedResult: []int{1, 2, 3},
			expectedError:  nil,
		},
		{
			name:           "Empty route index array",
			routeSeqArray:  []string{},
			expectedResult: nil,
			expectedError:  nil,
		},
		{
			name:           "Invalid route index array",
			routeSeqArray:  []string{"1", "2", "invalid"},
			expectedResult: nil,
			expectedError:  errors.New(`route_index should be an array of numbers - strconv.Atoi: parsing "invalid": invalid syntax`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GetRouteIndexFromParam(tc.routeSeqArray)
			if !IntSliceEqual(result, tc.expectedResult) {
				t.Errorf("GetRouteIndexFromParam() failed for test case '%s'. Unexpected result. Expected: %v, Got: %v", tc.name, tc.expectedResult, result)
			}
			if (err == nil && tc.expectedError != nil) || (err != nil && tc.expectedError == nil) || (err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error()) {
				t.Errorf("GetRouteIndexFromParam() failed for test case '%s'. Unexpected error. Expected: %v, Got: %v", tc.name, tc.expectedError, err)
			}
		})
	}
}

func TestValidateSortVerb(t *testing.T) {
	testCases := []struct {
		name           string
		sort           string
		allowedValues  []string
		expectedResult string
		expectedError  error
	}{
		{
			name:           "Empty sort parameter",
			sort:           "",
			allowedValues:  []string{"field1", "field2", "field3"},
			expectedResult: "created_at DESC",
			expectedError:  nil,
		},
		{
			name:           "Valid sort parameter",
			sort:           "field1:desc",
			allowedValues:  []string{"field1", "field2", "field3"},
			expectedResult: "field1 DESC",
			expectedError:  nil,
		},
		{
			name:           "Sort parameter with unknown field",
			sort:           "unknown:asc",
			allowedValues:  []string{"field1", "field2", "field3"},
			expectedResult: "",
			expectedError:  errors.New("Unknown field in sort query parameter"),
		},
		{
			name:           "Invalid sort parameter format",
			sort:           "field1~invalid",
			allowedValues:  []string{"field1", "field2", "field3"},
			expectedResult: "",
			expectedError:  errors.New("sort should be in this format field.orderdirection"),
		},
		{
			name:           "Invalid sort parameter order",
			sort:           "field1:invalid",
			allowedValues:  []string{"field1", "field2", "field3"},
			expectedResult: "",
			expectedError:  errors.New("Order direction should be asc or desc"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ValidateSortVerb(tc.sort, tc.allowedValues)
			if result != tc.expectedResult {
				t.Errorf("ValidateSortVerb() failed for test case '%s'. Unexpected result. Expected: %s, Got: %s", tc.name, tc.expectedResult, result)
			}
			if (err == nil && tc.expectedError != nil) || (err != nil && tc.expectedError == nil) || (err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error()) {
				t.Errorf("ValidateSortVerb() failed for test case '%s'. Unexpected error. Expected: %v, Got: %v", tc.name, tc.expectedError, err)
			}
		})
	}
}

func IntSliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
