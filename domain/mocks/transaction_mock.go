package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Transaction struct {
	mock.Mock
}

func (m *Transaction) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	if args.Get(0) != nil {
		return args.Error(0)
	}

	return fn(ctx)
}
