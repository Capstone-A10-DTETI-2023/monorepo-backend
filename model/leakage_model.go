package model

import "gorm.io/gorm"

type NodePressureRef struct {
	NodeID    uint `gorm:"not null" json:"node_id"`
	Pressure  float64 `gorm:"not null" json:"pressure"`
	Node		Node   `gorm:"foreignKey:NodeID"`
	gorm.Model
}

func MigrateNodePressureRef(db *gorm.DB) {
	db.AutoMigrate(&NodePressureRef{})
}

func (n *NodePressureRef) TableName() string {
	return "syssetting_node_pressure_ref"
}

func (n *NodePressureRef) CreateNodePressureRef(db *gorm.DB) error {
	return db.Create(n).Error
}

func (n *NodePressureRef) UpdateNodePressureRef(db *gorm.DB) error {
	return db.Save(n).Error
}

func (n *NodePressureRef) DeleteNodePressureRef(db *gorm.DB) error {
	return db.Unscoped().Delete(n).Error
}
