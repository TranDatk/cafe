package domain

import (
	"context"
	"time"
)

const (
	TableUser = "users"
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
	Fetch(c context.Context) ([]User, error)
	GetByEmail(c context.Context, email string) (User, error)
	GetByID(c context.Context, id string) (User, error)
	AssignRole(c context.Context, user *User, role *Role) error
}
