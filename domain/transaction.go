package domain

import (
	"context"
)

type Transaction interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

func WithinTransactionResult[T any](t Transaction, ctx context.Context, fn func(ctx context.Context) (T, error)) (T, error) {
	var result T
	err := t.WithinTransaction(ctx, func(txCtx context.Context) error {
		var err error
		result, err = fn(txCtx)
		return err
	})
	return result, err
}
