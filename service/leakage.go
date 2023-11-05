package service

import (
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
	"gonum.org/v1/gonum/mat"
	"gorm.io/gorm"
)

type Error struct {
	Message string
}

func (e Error) Error() string {
	return e.Message
}

func CalculateSensMatrix(db *gorm.DB) (*mat.Dense, error) {
	var nodeCount int64
	if err := db.Model(&model.Node{}).Where("calc_leakage = ?", true).Count(&nodeCount).Error; err != nil {
		return nil, err
	}

	if nodeCount < 2 {
		return nil, Error{"Not enough nodes to calculate sensitivity matrix"}
	}

	var defLeakSens, defNonLeakSens float64
	rows, err := db.Table("syssettings_general").Select("syssettings_general.def_leakage_sensitivity, syssettings_general.def_non_leak_sensitivity").Rows()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rows.Scan(&defLeakSens, &defNonLeakSens)
	}

	var (
		leakSens []float64
		refPres []float64
	)
	rows, err = db.Table("nodes").Select("nodes.leakage_sens, syssetting_node_pressure_ref.pressure").Joins("left join syssetting_node_pressure_ref on syssetting_node_pressure_ref.node_id = nodes.id").Rows()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var refPre float64
		var leakSen float64
		rows.Scan(&leakSen, &refPre)
		if refPre == -1 {
			return nil, Error{"Reference pressure not set"}
		}
		refPres = append(refPres, refPre)
		leakSens = append(leakSens, leakSen)
	}


	sensMat := mat.NewDense(int(nodeCount), int(nodeCount), nil)

	for i := 0; i < int(nodeCount); i++ {
		for j := 0; j < int(nodeCount); j++ {
			leakSens := leakSens[i]
			if leakSens == -1 {
				leakSens = defLeakSens
			}
			
			var leakPres float64
			if i == j {
				leakPres = refPres[j] * leakSens - refPres[j]
			} else {
				leakPres = refPres[j] - refPres[j] * leakSens * defNonLeakSens
			}

			sensMat.Set(i, j, leakPres)
		}
	}

	return sensMat, nil
}

