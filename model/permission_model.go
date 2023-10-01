package model

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	Role_ID		uint   `gorm:"not null"`
	Read_Realtime_Data bool `gorm:"not null"`
	Read_Historical_Data bool `gorm:"not null"`
	Change_Actuator bool `gorm:"not null"`
}

func MigratePermission(db *gorm.DB) {
	db.AutoMigrate(&Permission{})
	db.FirstOrCreate(&Permission{}, Permission{
		Role_ID: 1,
		Read_Realtime_Data: true,
		Read_Historical_Data: true,
		Change_Actuator: true,
	})
}

func (p *Permission) TableName() string {
	return "permissions"
}

func (p *Permission) CreatePermission(db *gorm.DB) error {
	return db.Create(p).Error
}

func (p *Permission) UpdatePermission(db *gorm.DB) error {
	return db.Save(p).Error
}

func (p *Permission) DeletePermission(db *gorm.DB) error {
	return db.Delete(p).Error
}
