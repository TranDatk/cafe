package seed

import (
	"log"

	"cafe/bootstrap"
	"cafe/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Seed(db *gorm.DB, env *bootstrap.Env) {
	log.Println("Seeding data...")

	// 1. Create Roles
	roles := []domain.Role{
		{
			ID:          uuid.Must(uuid.NewV7()).String(),
			Name:        "Administrator",
			Slug:        domain.AdminRole,
			Description: "System administrator with full access",
		},
		{
			ID:          uuid.Must(uuid.NewV7()).String(),
			Name:        "User",
			Slug:        domain.UserRole,
			Description: "Regular user with limited access",
		},
	}

	if err := db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "slug"}}, DoNothing: true}).Create(&roles).Error; err != nil {
		log.Fatalf("failed to seed roles: %v", err)
	}

	// Fetch actual IDs (in case they already existed)
	var dbRoles []domain.Role
	db.Where("slug IN ?", []string{domain.AdminRole, domain.UserRole}).Find(&dbRoles)
	roleMap := make(map[string]string)
	for _, r := range dbRoles {
		roleMap[r.Slug] = r.ID
	}
	adminRoleID := roleMap[domain.AdminRole]

	// 2. Create Permissions
	permissions := []domain.Permission{
		{ID: uuid.Must(uuid.NewV7()).String(), Name: "Create User", Slug: "user.create"},
		{ID: uuid.Must(uuid.NewV7()).String(), Name: "Update User", Slug: "user.update"},
		{ID: uuid.Must(uuid.NewV7()).String(), Name: "Delete User", Slug: "user.delete"},
		{ID: uuid.Must(uuid.NewV7()).String(), Name: "Read User", Slug: "user.read"},
		{ID: uuid.Must(uuid.NewV7()).String(), Name: "List Users", Slug: "user.list"},
	}

	if err := db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "slug"}}, DoNothing: true}).Create(&permissions).Error; err != nil {
		log.Fatalf("failed to seed permissions: %v", err)
	}

	// Fetch actual IDs for permissions
	var dbPermissions []domain.Permission
	permissionSlugs := make([]string, len(permissions))
	for i, p := range permissions {
		permissionSlugs[i] = p.Slug
	}
	db.Where("slug IN ?", permissionSlugs).Find(&dbPermissions)

	// Batch Link Permissions to Admin Role
	var rolePermissions []map[string]interface{}
	for _, p := range dbPermissions {
		rolePermissions = append(rolePermissions, map[string]interface{}{
			"role_id":       adminRoleID,
			"permission_id": p.ID,
		})
	}
	if len(rolePermissions) > 0 {
		if err := db.Table("role_permissions").Clauses(clause.OnConflict{DoNothing: true}).Create(&rolePermissions).Error; err != nil {
			log.Fatalf("failed to link permissions to admin role: %v", err)
		}
	}

	// 3. Create Admin User
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(env.SeedAdminPassword), bcrypt.DefaultCost)
	adminUser := domain.User{
		ID:       uuid.Must(uuid.NewV7()).String(),
		Name:     env.SeedAdminUsername,
		Email:    env.SeedAdminEmail,
		Password: string(hashedPassword),
	}

	if err := db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "email"}}, DoNothing: true}).Create(&adminUser).Error; err != nil {
		log.Fatalf("failed to seed admin user: %v", err)
	}

	// Fetch actual ID for admin user
	db.Where("email = ?", env.SeedAdminEmail).First(&adminUser)

	// 4. Assign Admin Role to User
	userRoleLink := map[string]interface{}{
		"user_id": adminUser.ID,
		"role_id": adminRoleID,
	}
	if err := db.Table("user_roles").Clauses(clause.OnConflict{DoNothing: true}).Create(&userRoleLink).Error; err != nil {
		log.Fatalf("failed to assign admin role to user: %v", err)
	}

	log.Println("Seeding completed successfully.")
}
