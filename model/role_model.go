package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name 		string `gorm:"not null"`
}

func MigrateRole(db *gorm.DB) {
	db.AutoMigrate(&Role{})
	db.FirstOrCreate(&Role{}, Role{
		Name: "SUPERADMIN",
	})
}

func (r *Role) TableName() string {
	return "roles"
}

func (r *Role) CreateRole(db *gorm.DB) error {
	return db.Create(r).Error
}

func (r *Role) UpdateRole(db *gorm.DB) error {
	return db.Save(r).Error
}

func (r *Role) DeleteRole(db *gorm.DB) error {
	return db.Delete(r).Error
}
