package repository

import (
	"context"

	"cafe/domain"
	"cafe/internal/pagination_util"

	"gorm.io/gorm"
)

var userFilterWhitelist = map[string]string{
	"id":        "id",
	"name":      "name",
	"email":     "email",
	"createdAt": "created_at",
	"updatedAt": "updated_at",
}

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

func (ur *userRepository) Fetch(c context.Context, opts domain.UserFetchOptions) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	db := ur.getDB(c).WithContext(c).Table(ur.table)

	db = db.Scopes(pagination_util.Filter(opts.Filters, userFilterWhitelist))

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Scopes(
		pagination_util.Sort(opts.Field, opts.SortOrder, userFilterWhitelist, userFilterWhitelist["id"]),
		pagination_util.Paginate(opts.Page, opts.PageSize),
	).Find(&users).Error

	return users, total, err
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
