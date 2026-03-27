package domain

const (
	TablePermission = "permissions"
)

type Permission struct {
	ID   string `gorm:"primaryKey" json:"id"`
	Name string `gorm:"column:name" json:"name"`
	Slug string `gorm:"column:slug;unique" json:"slug"`
}
