package domain

import "context"

type Transaction interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
