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

type RoleRepository interface {
	Create(c context.Context, role *Role) error
	Fetch(c context.Context) ([]Role, error)
	GetByID(c context.Context, id string) (Role, error)
	GetByName(c context.Context, name string) (Role, error)
}
