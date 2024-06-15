package utils

import "strconv"

func ConvertSliceOfStringToSliceOfInt(s []string) ([]int, error) {
	result := make([]int, len(s))
	for i, str := range s {
		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		result[i] = num
	}
	return result, nil
}

func ConvertSliceOfStringToSliceOfInt32(s []string) ([]int32, error) {
	result := make([]int32, len(s))
	for i, str := range s {
		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		result[i] = int32(num)
	}
	return result, nil
}

func CheckValueExistsInIntSlice(inputVal int32, intSlice []int32) bool {
	if intSlice == nil || len(intSlice) == 0 {
		return false
	}

	for _, val := range intSlice {
		if val == inputVal {
			return true
		}
	}
	return false
}

func ConvertSliceOfStringToSliceOfInt64(s []string) ([]int64, error) {
	result := make([]int64, len(s))
	for i, str := range s {
		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		result[i] = int64(num)
	}
	return result, nil
}
