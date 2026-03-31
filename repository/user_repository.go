package repository

import (
	"context"

	"cafe/domain"

	"gorm.io/gorm"
)

type userRepository struct {
	database *gorm.DB
	table    string
}

func NewUserRepository(db *gorm.DB, table string) domain.UserRepository {
	return &userRepository{
		database: db,
		table:    table,
	}
}

func (ur *userRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(TransactionKey).(*gorm.DB); ok {
		return tx
	}
	return ur.database
}

func (ur *userRepository) Create(c context.Context, user *domain.User) error {
	return ur.getDB(c).WithContext(c).Table(ur.table).Create(user).Error
}

func (ur *userRepository) Fetch(c context.Context) ([]domain.User, error) {
	var users []domain.User
	err := ur.getDB(c).WithContext(c).Table(ur.table).Find(&users).Error
	return users, err
}

func (ur *userRepository) GetByEmail(c context.Context, email string) (domain.User, error) {
	var user domain.User
	err := ur.getDB(c).WithContext(c).Table(ur.table).Where("email = ?", email).First(&user).Error
	return user, err
}

func (ur *userRepository) GetByID(c context.Context, id string) (domain.User, error) {
	var user domain.User
	err := ur.getDB(c).WithContext(c).Table(ur.table).Where("id = ?", id).First(&user).Error
	return user, err
}

func (ur *userRepository) AssignRole(c context.Context, user *domain.User, role *domain.Role) error {
	return ur.getDB(c).WithContext(c).Model(user).Association(domain.UserAssociationRoles).Append(role)
}
