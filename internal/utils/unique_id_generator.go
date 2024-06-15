package utils

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/shahbaz275817/prismo/pkg/logger"
)

var UUIDGenerator = uuid.NewRandom

func GenUUID() string {
	uid, err := UUIDGenerator()
	if err != nil {
		logger.WithContext(context.Background()).Errorf("Failure while generating uuid %s", err.Error())
		return "erroruid-xxxx-xxxx-xxxx-xx" + strconv.FormatInt(time.Now().Unix(), 10)
	}
	return uid.String()
}
