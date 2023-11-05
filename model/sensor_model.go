package model

import "gorm.io/gorm"

type Error struct {
	Message string
}

func (e Error) Error() string {
	return e.Message
}

type Sensor struct {
	gorm.Model
	Name 		string `gorm:"not null"`
	SensorType 	int    `gorm:"not null; default:0"`
	Unit 		string `gorm:"not null"`
	Interval 	int    `gorm:"not null; default:60"`
	Tolerance 	int    `gorm:"not null; default:10"`
	Alarm 		bool   `gorm:"not null; default:false"`
	AlarmType 	int
	AlarmLow 	float64
	AlarmHigh 	float64
	Node		Node   `gorm:"foreignKey:NodeID"`
	NodeID		uint
}

var SensorType = map[int]string{
	0: "Generic",
	1: "Pressure",
	2: "Flow",
	3: "Temperature",
}

var AlarmType = map[int]string{
	0: "Tidak ada alarm",
	1: "Alarm ketika nilai sensor lebih rendah dari nilai ambang bawah",
	2: "Alarm ketika nilai sensor lebih tinggi dari nilai ambang atas",
	3: "Alarm ketika nilai sensor berada di luar nilai ambang bawah dan atas",
}

func MigrateSensor(db *gorm.DB) {
	db.AutoMigrate(&Sensor{})
}

func (s *Sensor) TableName() string {
	return "sensors"
}

func (s *Sensor) CreateSensor(db *gorm.DB) error {
	if s.SensorType < 0 || s.SensorType > 3 {
		return Error{"Tipe sensor tidak valid"}
	}
	if s.Alarm && (s.AlarmType < 0 || s.AlarmType > 3) {
		return Error{"Tipe alarm tidak valid"}
	}
	if s.AlarmHigh < s.AlarmLow {
		return Error{"Nilai ambang atas harus lebih besar dari nilai ambang bawah"}
	}
	return db.Create(s).Error
}

func (s *Sensor) UpdateSensor(db *gorm.DB) error {
	return db.Save(s).Error
}

func (s *Sensor) DeleteSensor(db *gorm.DB) error {
	return db.Unscoped().Delete(s).Error
}

func (s *Sensor) GetSensor(db *gorm.DB) error {
	return db.First(s).Error
}