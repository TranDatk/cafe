package repository

import (
	"cafe/domain"
	"context"

	"gorm.io/gorm"
)

type contextKey string

const (
	TransactionKey contextKey = "transaction"
)

type transaction struct {
	db *gorm.DB
}

func NewTransaction(db *gorm.DB) domain.Transaction {
	return &transaction{
		db: db,
	}
}

func (t *transaction) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, TransactionKey, tx))
	})
}
