package repository

import (
	"context"

	"cafe/domain"
	"cafe/internal/pagination_util"

	"gorm.io/gorm"
)

type roleRepository struct {
	database *gorm.DB
	table    string
}

var roleFilterWhitelist = map[string]string{
	"id":   "id",
	"name": "name",
	"slug": "slug",
}

func NewRoleRepository(db *gorm.DB, table string) domain.RoleRepository {
	return &roleRepository{
		database: db,
		table:    table,
	}
}

func (rr *roleRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(TransactionKey).(*gorm.DB); ok {
		return tx
	}
	return rr.database
}

func (rr *roleRepository) Create(c context.Context, role *domain.Role) error {
	return rr.getDB(c).WithContext(c).Table(rr.table).Create(role).Error
}

func (rr *roleRepository) Fetch(c context.Context, opts domain.RoleFetchOptions) ([]domain.Role, int64, error) {
	var roles []domain.Role
	var total int64

	db := rr.getDB(c).WithContext(c).Table(rr.table)

	db = db.Scopes(pagination_util.Filter(opts.Filters, roleFilterWhitelist))

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Scopes(
		pagination_util.Sort(opts.Field, opts.SortOrder, roleFilterWhitelist, roleFilterWhitelist["id"]),
		pagination_util.Paginate(opts.Page, opts.PageSize),
	).Find(&roles).Error

	return roles, total, err
}

func (rr *roleRepository) GetByID(c context.Context, id string) (domain.Role, error) {
	var role domain.Role
	err := rr.getDB(c).WithContext(c).Table(rr.table).Where("id = ?", id).First(&role).Error
	return role, err
}

func (rr *roleRepository) GetByName(c context.Context, name string) (domain.Role, error) {
	var role domain.Role
	err := rr.getDB(c).WithContext(c).Table(rr.table).Where("name = ?", name).First(&role).Error
	return role, err
}
