package model

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	Email		bool  	`gorm:"not null"`
	Whatsapp	bool  	`gorm:"not null"`
	Firebase	bool  	`gorm:"not null"`
	User		User 	`gorm:"unique; foreignKey:UserID"`
	UserID		uint
}

type WhatsAppNotificationRequest struct {
	UserID uint `json:"UserID"`
	Message string `json:"Message"`
	Schedule string `json:"Schedule"`
}

func MigrateNotification(db *gorm.DB) {
	db.AutoMigrate(&Notification{})
	db.FirstOrCreate(&Notification{}, Notification{
		UserID: 1,
		Email: true,
		Whatsapp: true,
		Firebase: true,
	})
}

func (n *Notification) TableName() string {
	return "notifications"
}

func (n *Notification) CreateNotification(db *gorm.DB) error {
	return db.Create(n).Error
}

func (n *Notification) UpdateNotification(db *gorm.DB) error {
	return db.Save(n).Error
}

func (n *Notification) DeleteNotification(db *gorm.DB) error {
	return db.Delete(n).Error
}
