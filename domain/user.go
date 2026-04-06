package domain

import (
	"context"
	"time"
)

const (
	TableUser            = "users"
	UserAssociationRoles = "Roles"
)

type User struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Email     string    `gorm:"column:email;unique" json:"email"`
	Password  string    `gorm:"column:password" json:"password"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Roles     []Role    `gorm:"many2many:user_roles" json:"roles"`
}

type UserRepository interface {
	Create(c context.Context, user *User) error
	Fetch(c context.Context, opts UserFetchOptions) ([]User, int64, error)
	GetByEmail(c context.Context, email string) (User, error)
	GetByID(c context.Context, id string) (User, error)
	AssignRole(c context.Context, user *User, role *Role) error
}

type UserFetchOptions struct {
	Page     int `form:"page" binding:"omitempty"`
	PageSize int `form:"page_size" binding:"omitempty"`
	SortOption
	Filters []Criterion
}

func (o *UserFetchOptions) Prepare() {
	if o.Page <= 0 {
		o.Page = 1
	}
	if o.PageSize <= 0 {
		o.PageSize = 10
	}
}

type UserUsecase interface {
	Fetch(c context.Context, opts UserFetchOptions) ([]User, int64, error)
}
