package locks

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"

	"github.com/shahbaz275817/prismo/pkg/cache/mocks"
	"github.com/shahbaz275817/prismo/pkg/errors"
)

func TestAtomicLockExecute_Lock_Success(t *testing.T) {
	client := &mocks.MockCacheClient{}
	config := map[KeyType]LockConfig{
		Def: {
			LockExpiry:    time.Duration(10) * time.Millisecond,
			RetryAttempts: uint(1),
			RetryDelay:    time.Duration(1) * time.Millisecond,
		},
	}

	atomicLock := NewAtomicLock(client, config)

	tests := []struct {
		name       string
		key        string
		keyType    KeyType
		atomicFunc func(lockState *LockState) ([]interface{}, error)
		mockFunc   func()
		state      *LockState
		wantErr    bool
	}{
		{
			name:    "test successful lock on a default key",
			key:     "abc",
			keyType: Def,
			atomicFunc: func(lockState *LockState) ([]interface{}, error) {
				return nil, nil
			},
			mockFunc: func() {
				client.On("SetNX", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redis.NewBoolResult(true, nil)).Once()
				client.On("Del", mock.Anything, mock.Anything).Return(redis.NewIntResult(1, nil)).Once()
			},
			state:   &LockState{Mutex: sync.Mutex{}},
			wantErr: false,
		},
		{
			name:    "test different KeyType config",
			key:     "xyz",
			keyType: Def,
			atomicFunc: func(lockState *LockState) ([]interface{}, error) {
				return nil, nil
			},
			mockFunc: func() {
				client.On("SetNX", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redis.NewBoolResult(true, nil)).Once()
				client.On("Del", mock.Anything, mock.Anything).Return(redis.NewIntResult(1, nil)).Once()
			},
			state:   &LockState{Mutex: sync.Mutex{}},
			wantErr: false,
		},
		{
			name:    "test reentrant lock when key is already present in lock state",
			key:     "abc",
			keyType: Def,
			atomicFunc: func(lockState *LockState) ([]interface{}, error) {
				return nil, nil
			},
			mockFunc: func() {
				// No SetNX or Del calls expected because it should recognize the reentrant lock
			},
			state: &LockState{
				Mutex:    sync.Mutex{},
				LockKeys: map[string]struct{}{"abc": {}},
			},
			wantErr: false,
		}, {
			name:    "test re-entrant lock when atomic lock is called inside inner func",
			key:     "xyz",
			keyType: Def,
			atomicFunc: func(lockState *LockState) ([]interface{}, error) {
				innerFunc := func(innerLockState *LockState) ([]interface{}, error) {
					// Do something here if needed
					return nil, nil
				}
				ctx := context.Background()
				_, err := atomicLock.Execute(ctx, "xyz", Def, innerFunc, lockState)
				if err != nil {
					return nil, err
				}
				return nil, nil
			},
			mockFunc: func() {
				client.On("SetNX", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redis.NewBoolResult(true, nil)).Once()
				client.On("Del", mock.Anything, mock.Anything).Return(redis.NewIntResult(1, nil)).Once()
			},
			state:   &LockState{Mutex: sync.Mutex{}},
			wantErr: false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			// Clear expectations from previous runs
			client.ExpectedCalls = nil
			client.Calls = nil

			//Execute mock function
			tt.mockFunc()

			//Execute the main testing function
			_, err := atomicLock.Execute(ctx, tt.key, tt.keyType, tt.atomicFunc, tt.state)

			if (err != nil) != tt.wantErr {
				t.Errorf("Test Failed: got error = %v, wantErr %v", err, tt.wantErr)
			}

			// Assert that all expected calls were made
			client.AssertExpectations(t)

		})

	}
}

func TestAtomicLock_Execute_Lock_Failures(t *testing.T) {
	client := &mocks.MockCacheClient{}
	config := map[KeyType]LockConfig{
		Def: {
			LockExpiry:    time.Duration(10) * time.Millisecond,
			RetryAttempts: uint(1),
			RetryDelay:    time.Duration(1) * time.Millisecond,
		},
	}

	atomicLock := NewAtomicLock(client, config)

	tests := []struct {
		name        string
		key         string
		keyType     KeyType
		atomicFunc  func(lockState *LockState) ([]interface{}, error)
		mockFunc    func()
		state       *LockState
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "when lock acquisition fails with no retries",
			key:     "def",
			keyType: Def,
			atomicFunc: func(lockState *LockState) ([]interface{}, error) {
				return nil, nil
			},
			mockFunc: func() {
				client.On("SetNX", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redis.NewBoolResult(false, nil)).Once()
			},
			state:       &LockState{Mutex: sync.Mutex{}},
			wantErr:     true,
			expectedErr: errors.EntityLockedError{},
		},

		{
			name:    "when lock is acquired on second attempt",
			key:     "abc",
			keyType: Def,
			atomicFunc: func(lockState *LockState) ([]interface{}, error) {
				return nil, nil
			},
			mockFunc: func() {
				client.On("SetNX", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redis.NewBoolResult(false, nil)).Once()
				client.On("SetNX", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redis.NewBoolResult(true, nil)).Once()
				client.On("Del", mock.Anything, mock.Anything).Return(redis.NewIntResult(1, nil)).Once()
			},
			state: &LockState{
				Mutex: sync.Mutex{},
			},
			wantErr: false,
		},
		{
			name:    "when client keeps throwing error while retry attempts are exhausted",
			key:     "abc",
			keyType: Def,
			atomicFunc: func(lockState *LockState) ([]interface{}, error) {
				return nil, nil
			},
			mockFunc: func() {
				client.On("SetNX", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redis.NewBoolResult(false, errors.New("Lock Error"))).Once()
			},
			state: &LockState{
				Mutex: sync.Mutex{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			// Clear expectations from previous runs
			client.ExpectedCalls = nil
			client.Calls = nil

			//Execute mock function
			tt.mockFunc()

			//Execute the main testing function
			_, err := atomicLock.Execute(ctx, tt.key, tt.keyType, tt.atomicFunc, tt.state)

			if (err != nil) != tt.wantErr {
				t.Errorf("Test Failed: got error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.expectedErr != nil {
				assert.ErrorAs(t, err, &tt.expectedErr)
			}

			// Assert that all expected calls were made
			client.AssertExpectations(t)

		})

	}

}
