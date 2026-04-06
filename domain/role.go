package domain

import (
	"context"
	"time"
)

const (
	TableRole = "roles"
	AdminRole = "admin"
	UserRole  = "user"
)

type Role struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`
	Slug        string    `gorm:"column:slug;unique" json:"slug"`
	Description string    `gorm:"column:description" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
}

type RoleFetchOptions struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	SortOption
	Filters []Criterion
}

type RoleUsecase interface {
	Fetch(c context.Context, opts RoleFetchOptions) ([]Role, int64, error)
}

type RoleRepository interface {
	Create(c context.Context, role *Role) error
	Fetch(c context.Context, opts RoleFetchOptions) ([]Role, int64, error)
	GetByID(c context.Context, id string) (Role, error)
	GetByName(c context.Context, name string) (Role, error)
}
