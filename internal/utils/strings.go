package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func ToFloat64(convertibleString string) (float64, error) {
	return strconv.ParseFloat(convertibleString, 64)
}

func Float64ToString(number float64) string {
	return strconv.FormatFloat(number, 'f', -1, 64)
}

func ToString(integer int64) string {
	return strconv.FormatInt(integer, 10)
}

func InterfaceToString(val interface{}) string {
	return InterfaceToStringWithDefault(val, "")
}

func InterfaceToStringWithDefault(val interface{}, def string) string {
	str := def
	if val != nil && reflect.TypeOf(val).String() == "string" {
		str = val.(string)
	}
	return str
}

func IsEmptyStringPointer(val *string) bool {
	return val == nil || *val == ""
}

// GenerateStoreClient Tried to accommodate as many strings, but golang doesn't support negative look ahead but ruby does and code has it
func GenerateStoreClient(str string) string {
	str = strings.TrimSpace(str)
	reg1, _ := regexp.Compile(`\s+`)
	str = reg1.ReplaceAllString(str, "_")
	str = ToUnderScore(str, '_')
	str = strings.ReplaceAll(str, "_", "-")
	str = strings.ToLower(str)
	return str + "-engine"
}

// ToUnderScore Please refer this https://github.com/iancoleman/strcase
// Done some modifications to accommodate the change required
func ToUnderScore(s string, delimiter uint8) string {
	s = strings.Trim(s, " ")
	n := ""
	for i, v := range s {
		// treat acronyms as words, eg for JSONData -> JSON is a whole word
		nextCaseIsChanged := false
		if i+1 < len(s) {
			next := s[i+1]
			vIsCap := v >= 'A' && v <= 'Z'
			vIsLow := v >= 'a' && v <= 'z'
			nextIsCap := next >= 'A' && next <= 'Z'
			nextIsLow := next >= 'a' && next <= 'z'
			if (vIsCap && nextIsLow) || (vIsLow && nextIsCap) {
				nextCaseIsChanged = true
			}
			if i-1 >= 0 && checkIgnore(s[i-1]) && nextCaseIsChanged {
				nextCaseIsChanged = false
			}
		}

		if i > 0 && n[len(n)-1] != delimiter && nextCaseIsChanged {
			// add underscore if next letter case type is changed
			if v >= 'A' && v <= 'Z' {
				n += string(delimiter) + string(v)
			} else if v >= 'a' && v <= 'z' {
				n += string(v) + string(delimiter)
			}
		} else if v == ' ' || v == '_' || v == '-' {
			// replace spaces/underscores with delimiters
			if checkIgnore(uint8(v)) {
				n += string(v)
			} else {
				n += string(delimiter)
			}
		} else {
			n = n + string(v)
		}
	}
	return n
}

func checkIgnore(u uint8) bool {
	return !(u >= 'A' && u <= 'Z' || u >= 'a' && u <= 'z')
}

func StringInSlice(strSlice []string, s string) bool {
	for _, v := range strSlice {
		if v == s {
			return true
		}
	}
	return false
}

func CompareStringSlices(a, b []string) bool {
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

func IsAnyStringInSlice(stringsToFind, searchSlice []string) bool {
	elementsInSearchSlice := make(map[string]bool)

	for _, item := range searchSlice {
		elementsInSearchSlice[item] = true
	}

	for _, item := range stringsToFind {
		if elementsInSearchSlice[item] {
			return true
		}
	}

	return false
}

func TrimStrings(arr []string) []string {
	if arr == nil {
		return nil
	}
	trimmedArr := make([]string, len(arr))
	for i, s := range arr {
		trimmedArr[i] = strings.TrimSpace(s)
	}
	return trimmedArr
}

func ConvertSliceOfIntToSliceOfString(s []int) []string {
	result := make([]string, len(s))
	for i, num := range s {
		result[i] = strconv.Itoa(num)
	}
	return result
}

func ConvertSliceOfInt32ToSliceOfString(s []int32) []string {
	result := make([]string, len(s))
	for i, num := range s {
		result[i] = strconv.Itoa(int(num))
	}
	return result
}

func CovertErrorMessagesToString(errs []error) string {
	errStrings := make([]string, len(errs))
	for i, err := range errs {
		errStrings[i] = err.Error()
	}
	return strings.Join(errStrings, "\n")
}

func BuildLockKey(arguments ...interface{}) string {
	key := "lock"
	for _, arg := range arguments {
		key = fmt.Sprintf("%s-%v", key, arg)
	}
	return key
}

func ConvertInterfaceToStringSlice(in interface{}) ([]string, error) {
	slice, ok := in.([]interface{})

	if !ok {
		return nil, fmt.Errorf("input is not a slice of interfaces: %T", in)
	}
	result := make([]string, len(slice))
	for i, v := range slice {
		if str, ok := v.(string); ok {
			result[i] = str
		} else {
			return nil, fmt.Errorf("value at index %d is not a string: %v", i, v)
		}
	}
	return result, nil
}
