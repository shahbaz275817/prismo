package config

import (
	"time"

	config2 "github.com/shahbaz275817/prismo/pkg/config"
	"github.com/shahbaz275817/prismo/pkg/locks"
)

func newAtomicLockConfig() map[locks.KeyType]locks.LockConfig {
	a := make(map[locks.KeyType]locks.LockConfig)
	a[locks.Def] = locks.LockConfig{
		LockExpiry:    time.Duration(config2.MustGetInt("AL_DEF_LOCK_EXPIRY_MS")) * time.Millisecond,
		RetryAttempts: uint(config2.MustGetInt("AL_DEF_RETRY_ATTEMPTS")),
		RetryDelay:    time.Duration(config2.MustGetInt("AL_DEF_RETRY_DELAY")) * time.Millisecond,
	}

	return a
}
