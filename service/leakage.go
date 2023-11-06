package service

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
	"github.com/jackc/pgx/v5"
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
		nonLeakSens []float64
		refPres []float64
	)
	rows, err = db.Table("nodes").Select("nodes.leakage_sens, nodes.non_leak_sens, syssetting_node_pressure_ref.pressure").Joins("left join syssetting_node_pressure_ref on syssetting_node_pressure_ref.node_id = nodes.id").Order("nodes.id ASC").Rows()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var refPre float64
		var leakSen float64
		var nonLeakSen float64
		rows.Scan(&leakSen, &nonLeakSen, &refPre)
		if refPre == -1 {
			return nil, Error{"Reference pressure not set"}
		}
		refPres = append(refPres, refPre)
		leakSens = append(leakSens, leakSen)
		nonLeakSens = append(nonLeakSens, nonLeakSen)
	}


	sensMat := mat.NewDense(int(nodeCount), int(nodeCount), nil)

	for i := 0; i < int(nodeCount); i++ {
		for j := 0; j < int(nodeCount); j++ {
			leakSens := leakSens[i]
			nonLeakSens := nonLeakSens[i]
			if leakSens == -1 {
				leakSens = defLeakSens
			}
			
			var leakPres float64
			if i == j {
				leakPres = refPres[j] * (1-leakSens) - refPres[j]
			} else {
				leakPres = refPres[j] * nonLeakSens - refPres[j]
			}

			sensMat.Set(i, j, leakPres)
		}
	}

	return sensMat, nil
}

func GetNodesId(db *gorm.DB) ([]int, int, error) {
	var nodesId []int
	rows, err := db.Table("nodes").Select("nodes.id").Where("calc_leakage = ?", true).Order("nodes.id ASC").Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var nodeId int
		rows.Scan(&nodeId)
		nodesId = append(nodesId, nodeId)
	}

	nodeCount := len(nodesId)

	return nodesId, nodeCount, nil
}

func GetRefPressure(db *gorm.DB) (map[int]float64, error) {
	refPresData := make(map[int]float64)

	rows, err := db.Table("nodes").Select("nodes.id, syssetting_node_pressure_ref.pressure").Joins("left join syssetting_node_pressure_ref on syssetting_node_pressure_ref.node_id = nodes.id").Order("nodes.id ASC").Rows()
	if err != nil {
		return refPresData, err
	}
	defer rows.Close()
	for rows.Next() {
		var nodeId int
		var refPre float64
		rows.Scan(&nodeId, &refPre)
		if refPre == -1 {
			return refPresData, Error{"Reference pressure not set"}
		}
		refPresData[nodeId] = refPre
	}
	return refPresData, nil
}

func GetLatestSensorData(db *gorm.DB, dbTs *pgx.Conn) (map[int]float64, error) {
	sensorPresData := make(map[int]float64)

	var sensorsId []int
	rows, err := db.Table("sensors").Select("sensors.id").Joins("left join nodes on nodes.id = sensors.node_id").Where("nodes.calc_leakage = ?", true).Where("sensors.sensor_type = ?", 1).Order("nodes.id ASC").Rows()
	if err != nil {
		return sensorPresData, err
	}
	defer rows.Close()
	for rows.Next() {
		var sensId int
		rows.Scan(&sensId)
		sensorsId = append(sensorsId, sensId)
	}

	for _, id := range sensorsId {
		sensorId := strconv.Itoa(id)
		query := fmt.Sprintf("SELECT value FROM %s WHERE sensor_id = '%s' ORDER BY timestamp DESC LIMIT 1", "sensor_data", sensorId)
		data := dbTs.QueryRow(context.Background(), query)
		value := "0"
		if err := data.Scan(&value); (err != nil && err != pgx.ErrNoRows) {
			return sensorPresData, err
		}
		valueFloat, _ := strconv.ParseFloat(value, 64)
		sensorPresData[id] = valueFloat
	}
	log.Println(sensorPresData)
	return sensorPresData, nil
}

func CalculateResidualMatrix(db *gorm.DB, dbTs *pgx.Conn) (*mat.Dense, error) {
	nodesId, nodeCount, err := GetNodesId(db)
	if err != nil {
		return nil, err
	}
	latestSensData, err := GetLatestSensorData(db, dbTs)
	if err != nil {
		return nil, err
	}
	refPresData, err := GetRefPressure(db)
	if err != nil {
		return nil, err
	}


	resMat := mat.NewDense(nodeCount, 1, nil)

	for i := 0; i < nodeCount; i++ {
		nodeId := nodesId[i]
		refPres := refPresData[nodeId]
		latestPres := latestSensData[nodeId]
		residual := latestPres - refPres

		resMat.Set(i, 0, residual)
	}

	return resMat, nil
}
