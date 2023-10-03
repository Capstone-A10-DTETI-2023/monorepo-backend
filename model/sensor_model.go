package model

import "gorm.io/gorm"

type Sensor struct {
	gorm.Model
	Name 		string `gorm:"not null"`
	Unit 		string `gorm:"not null"`
	Node		Node   `gorm:"foreignKey:NodeID"`
	NodeID		uint
}

func MigrateSensor(db *gorm.DB) {
	db.AutoMigrate(&Sensor{})
}

func (s *Sensor) TableName() string {
	return "sensors"
}

func (s *Sensor) CreateSensor(db *gorm.DB) error {
	return db.Create(s).Error
}

func (s *Sensor) UpdateSensor(db *gorm.DB) error {
	return db.Save(s).Error
}

func (s *Sensor) DeleteSensor(db *gorm.DB) error {
	return db.Delete(s).Error
}

func (s *Sensor) GetSensor(db *gorm.DB) error {
	return db.First(s).Error
}
