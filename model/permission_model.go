package model

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	Read_Realtime_Data 		bool `gorm:"not null default:true"`
	Read_Historical_Data 	bool `gorm:"not null default:false"`
	Change_Actuator 		bool `gorm:"not null default:false"`
	User_Management 		bool `gorm:"not null default:false"`
	Node_Management 		bool `gorm:"not null default:false"`
	Role 					Role `gorm:"unique; foreignKey:RoleID"`
	RoleID 					uint
}

func MigratePermission(db *gorm.DB) {
	db.AutoMigrate(&Permission{})
}

func BootstrapPermission(db *gorm.DB) {
	db.FirstOrCreate(&Permission{}, Permission{
		RoleID: 1,
		Read_Realtime_Data: true,
		Read_Historical_Data: true,
		Change_Actuator: true,
		User_Management: true,
		Node_Management: true,
	})
	db.FirstOrCreate(&Permission{}, Permission{
		RoleID: 2,
		Read_Realtime_Data: true,
		Read_Historical_Data: true,
		Change_Actuator: true,
		User_Management: false,
		Node_Management: false,
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
