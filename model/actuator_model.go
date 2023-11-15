package model

import "gorm.io/gorm"

type Actuator struct {
	gorm.Model
	Name 		string `gorm:"not null"`
	ActuatorType 	int    `gorm:"not null; default:0"`	// 0: ON_OFF, 1: Setpoint
	Unit 		string `gorm:"not null"`
	Interval 	int    `gorm:"not null; default:60"`
	Node		Node   `gorm:"foreignKey:NodeID"`
	NodeID		uint  `gorm:"not null"`
}

var ActuatorType = map[int]string{
	0: "ON_OFF",
	1: "Setpoint",
}

func MigrateActuator(db *gorm.DB) {
	db.AutoMigrate(&Actuator{})
}

func (a *Actuator) TableName() string {
	return "actuators"
}

func (a *Actuator) CreateActuator(db *gorm.DB) error {
	if a.ActuatorType < 0 || a.ActuatorType > 1 {
		return Error{"Tipe actuator tidak valid"}
	}
	return db.Create(a).Error
}

func (a *Actuator) UpdateActuator(db *gorm.DB) error {
	return db.Save(a).Error
}

func (a *Actuator) DeleteActuator(db *gorm.DB) error {
	return db.Unscoped().Delete(a).Error
}

func (a *Actuator) GetActuator(db *gorm.DB) error {
	return db.First(a).Error
}