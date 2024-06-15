package utils

import "fmt"

func CreateAWBLegTypeLockKey(awbNumber string, legType string) string {
	return fmt.Sprintf("%s%s", awbNumber, legType)
}
