package service

import (
	"log"

	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
	"gonum.org/v1/gonum/mat"
	"gorm.io/gorm"
)

func CalculateSensMatrix(db *gorm.DB) error {
	var nodeCount int64
	if err := db.Model(&model.Node{}).Count(&nodeCount).Error; err != nil {
		return err
	}

	sensMat := mat.NewDense(int(nodeCount), int(nodeCount), nil)
	log.Println(sensMat)
	return nil
}