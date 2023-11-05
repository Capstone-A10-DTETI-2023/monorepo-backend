package model

import "gorm.io/gorm"

type SystemSetting struct {
	gorm.Model
	WebsiteName			string `gorm:"not null"`
	WebsiteDescription	string `gorm:"not null"`
	DefLeakageSensitivity 		float64 `gorm:"not null"`
	DefNonLeakSensitivity 		float64 `gorm:"not null"`
	WhatsappCooldownSecs 	int `gorm:"not null"`
}

func MigrateSystemSetting(db *gorm.DB) {
	db.AutoMigrate(&SystemSetting{})
}

func BootstrapSystemSetting(db *gorm.DB) {
	db.FirstOrCreate(&SystemSetting{}, SystemSetting{
		WebsiteName: "Sistem Informasi Monitoring Tekanan Air",
		WebsiteDescription: "Sistem Informasi Monitoring Tekanan Air",
		DefLeakageSensitivity: 0.5,
		DefNonLeakSensitivity: 0.5,
		WhatsappCooldownSecs: 30,
	})
}

func (s *SystemSetting) TableName() string {
	return "syssettings_general"
}

func (s *SystemSetting) CreateSystemSetting(db *gorm.DB) error {
	return db.Create(s).Error
}

func (s *SystemSetting) UpdateSystemSetting(db *gorm.DB) error {
	return db.Save(s).Error
}

func (s *SystemSetting) DeleteSystemSetting(db *gorm.DB) error {
	return db.Unscoped().Delete(s).Error
}

func (s *SystemSetting) GetSystemSetting(db *gorm.DB) error {
	return db.First(s).Error
}
