package repository

import (
	"context"
)

type BaseRepository interface {
	Transact(ctx context.Context, fn func(context.Context) error) error
}

type BaseRepositoryImpl struct {
	DB Accessor
}

// Transact Wrapper to run a function in a transaction block. Supports nested transaction, the outermost transaction block
// will be committed or rolled back.
func (cr *BaseRepositoryImpl) Transact(ctx context.Context, fn func(context.Context) error) error {
	return cr.DB.Transact(ctx, fn)
}
