package controller

import (
	"fmt"

	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/service"
	"github.com/gofiber/fiber/v2"
	"gonum.org/v1/gonum/mat"
	"gorm.io/gorm"
)

type LeakageController struct {
	DB *gorm.DB
}

func NewLeakageController(db *gorm.DB) *LeakageController {
	return &LeakageController{
		DB: db,
	}
}

func (c *LeakageController) GetSensMat(ctx *fiber.Ctx) error {
	sensMat, err := service.CalculateSensMatrix(c.DB)
	if err != nil {
		return err
	}

	fSensMat := mat.Formatted(sensMat, mat.FormatMATLAB())
	fSensMatStr := fmt.Sprintf("%.3g", fSensMat)

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    fSensMatStr,
	})
}