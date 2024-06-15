package locks

import (
	"context"
	defaultErrors "errors"
	"fmt"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"

	"github.com/shahbaz275817/prismo/pkg/cache"
	"github.com/shahbaz275817/prismo/pkg/errors"
	"github.com/shahbaz275817/prismo/pkg/logger"
)

// KeyType represents the type of lock.
type KeyType string

const (
	Def KeyType = "DEFAULT"
)

// LockConfig holds configuration for a particular type of lock.
type LockConfig struct {
	LockExpiry    time.Duration
	RetryAttempts uint
	RetryDelay    time.Duration
}

// LockState holds the state of the lock.
type LockState struct {
	sync.Mutex
	LockKeys map[string]struct{}
}

// AtomicLock provides methods for distributed locking.
type AtomicLock struct {
	client  cache.Client
	configs map[KeyType]LockConfig
}

// AtomicLock's Singleton Instance
var atomicLockInstance *AtomicLock

// NewAtomicLock creates a new AtomicLock with the given cache client and configuration.
func NewAtomicLock(client cache.Client, config map[KeyType]LockConfig) *AtomicLock {
	return &AtomicLock{
		client:  client,
		configs: config,
	}
}

// Execute executes a function within a lock for a given key.
func (a *AtomicLock) Execute(ctx context.Context, key string, keyType KeyType, f func(*LockState) ([]interface{}, error), state *LockState) ([]interface{}, error) {
	config := a.getConfig(keyType)

	if !a.isReentrant(state, key) {
		err := a.getLockWithRetries(ctx, key, state, config.LockExpiry, config.RetryAttempts, config.RetryDelay)
		if err != nil {
			return nil, err
		}
		defer a.releaseLock(state, key)
	}

	return f(state)
}

// isReentrant checks if the current goroutine already owns the lock.
func (a *AtomicLock) isReentrant(state *LockState, key string) bool {
	state.Lock()
	defer state.Unlock()

	_, exists := state.LockKeys[key]
	return exists
}

// getConfig gets the lock configuration for the given key type, falling back to default if not found.
func (a *AtomicLock) getConfig(keyType KeyType) LockConfig {
	if config, ok := a.configs[keyType]; ok {
		return config
	}
	return a.configs["default"]
}

func (a *AtomicLock) getLockWithRetries(ctx context.Context, key string, state *LockState, expiry time.Duration, attempts uint, delay time.Duration) error {
	err := retry.Do(
		func() error {
			return a.lock(state, key, expiry)
		},
		retry.Attempts(attempts), retry.Delay(delay))
	if err != nil {
		return handleError(ctx, err)
	}
	return nil
}

// lock attempts to acquire a lock on the given key with the specified expiry.
func (a *AtomicLock) lock(state *LockState, key string, expiry time.Duration) error {
	err := a.setKeyInCache(key, expiry)
	if err != nil {
		return err
	}
	a.setKeyInLockState(state, key)
	return nil
}

// releaseLock releases the lock on the given key if owned.
func (a *AtomicLock) releaseLock(state *LockState, key string) error {
	_, err := a.deleteKeyInCache(key)
	if err != nil {
		return err
	}
	a.deleteKeyFromLockState(state, key)
	return nil
}

func (a *AtomicLock) setKeyInLockState(state *LockState, key string) {
	state.Lock()
	defer state.Unlock()

	if state.LockKeys == nil {
		state.LockKeys = make(map[string]struct{})
	}
	state.LockKeys[key] = struct{}{}
}

func (a *AtomicLock) deleteKeyFromLockState(l *LockState, key string) {
	l.Lock()
	defer l.Unlock()

	if l.LockKeys != nil {
		delete(l.LockKeys, key)
	}
}

func (a *AtomicLock) setKeyInCache(key string, expiry time.Duration) error {
	res, err := a.client.SetNX(context.Background(), key, true, expiry).Result()
	if err != nil {
		return err
	}
	if !res {
		return errors.NewEntityLockedError("unable_to_acquire_lock", &errors.ErrDetails{
			Message: fmt.Sprintf("Unable to acquire lock on key %s", key),
		})
	}
	return nil
}

func (a *AtomicLock) deleteKeyInCache(key string) (bool, error) {
	res, err := a.client.Del(context.Background(), key).Result()
	if err != nil {
		return false, err
	}
	return res != 0, err
}

func handleError(ctx context.Context, err error) error {
	isSafeCast := defaultErrors.As(err, &retry.Error{})
	if !isSafeCast {
		return err
	}
	hasInternalError := false
	// check all the errors to see if there is any other error type beside NewEntityLockedError
	for _, e := range err.(retry.Error) {
		if !defaultErrors.As(e, &errors.EntityLockedError{}) {
			hasInternalError = true
			break
		}
	}
	logger.WithContext(ctx).Warn(err.Error())
	if hasInternalError {
		return errors.NewInternalServerError("atomic_lock_error", &errors.ErrDetails{
			Message: fmt.Sprintf("Something went wrong while processing the request. Please try sometimes later."),
		})
	}
	// return only one NewEntityLockedError if all errors for each retry is a NewEntityLockedError
	return errors.NewEntityLockedError("unable_to_acquire_lock", &errors.ErrDetails{
		Message: "Unable to process your request at this moment. Please try sometimes later.",
	})
}
