package model

import "gorm.io/gorm"

type Node struct {
	gorm.Model
	Name 		string `gorm:"not null; unique"`
	Latitude	string `gorm:"not null"`
	Longitude	string `gorm:"not null"`
	CalcLeakage bool `gorm:"not null; default:false"`
	LeakageSens float64 `gorm:"not null; default:-1"`
}

func MigrateNode(db *gorm.DB) {
	db.AutoMigrate(&Node{})
}

func (n *Node) TableName() string {
	return "nodes"
}

func (n *Node) CreateNode(db *gorm.DB) error {
	return db.Create(n).Error
}

func (n *Node) UpdateNode(db *gorm.DB) error {
	return db.Save(n).Error
}

func (n *Node) DeleteNode(db *gorm.DB) error {
	return db.Unscoped().Delete(n).Error
}

func (n *Node) GetNode(db *gorm.DB) error {
	return db.First(n).Error
}
