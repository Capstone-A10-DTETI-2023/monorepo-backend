package service

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
	"github.com/jackc/pgx/v5"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
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
	sensorPresTS := make(map[int]time.Time)
	sensorTimeTolerance := make(map[int]int)

	var sensorsId []int

	rows, err := db.Table("sensors").Select("sensors.id, sensors.tolerance").Distinct("nodes.id").Joins("left join nodes on nodes.id = sensors.node_id").Where("nodes.calc_leakage = ?", true).Where("sensors.sensor_type = ?", 1).Order("nodes.id ASC").Rows()
	if err != nil {
		return sensorPresData, err
	}
	defer rows.Close()
	for rows.Next() {
		var sensId, sensTol int
		rows.Scan(&sensId, &sensTol)
		sensorsId = append(sensorsId, sensId)
		sensorTimeTolerance[sensId] = sensTol
	}

	for _, id := range sensorsId {
		sensorId := strconv.Itoa(id)
		query := fmt.Sprintf("SELECT timestamp,value FROM %s WHERE sensor_id = '%s' ORDER BY timestamp DESC LIMIT 1", "sensor_data", sensorId)
		data := dbTs.QueryRow(context.Background(), query)
		value := "0"
		var timestamp time.Time
		if err := data.Scan(&timestamp, &value); (err != nil && err != pgx.ErrNoRows) {
			return sensorPresData, err
		}
		valueFloat, _ := strconv.ParseFloat(value, 64)
		if err != nil {
			return sensorPresData, err
		}
		sensorPresData[id] = valueFloat
		sensorPresTS[id] = timestamp
	}

	for _, id := range sensorsId {
		if time.Since(sensorPresTS[id]).Seconds() > float64(sensorTimeTolerance[id]) {
			return sensorPresData, Error{"Sensor data is outdated. Check sensor connection."}
		}
	}
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

func GetLeakageNode(sensMat *mat.Dense, resMat *mat.Dense, db *gorm.DB) (int, error) {
	if sensMat == nil || resMat == nil {
		return -1, Error{"Sensitivity matrix or residual matrix is nil"}
	}

	if sensMat.RawMatrix().Rows != sensMat.RawMatrix().Cols {
		return -1, Error{"Sensitivity matrix is not square"}
	}

	if sensMat.RawMatrix().Cols != resMat.RawMatrix().Rows {
		return -1, Error{"Sensitivity matrix cols and residual matrix rows are not equal"}
	}

	nodesId, nodeCount, err := GetNodesId(db)
	if err != nil {
		return -1, err
	}

	var correlation []float64
	for i := 0; i < nodeCount; i++ {
		sensMatCol := mat.Col(nil, i, sensMat)
		// log.Println(sensMatCol)
		resMatCol := mat.Col(nil, 0, resMat)
		// log.Println(resMatCol)
		matSRStack := mat.NewDense(2, len(sensMatCol), nil)
		matSRStack.SetRow(0, sensMatCol)
		matSRStack.SetRow(1, resMatCol)
		// log.Printf("matSRStack: %.2g", mat.Formatted(matSRStack, mat.FormatMATLAB()))
		matRRStack := mat.NewDense(2, len(resMatCol), nil)
		matRRStack.SetRow(0, resMatCol)
		matRRStack.SetRow(1, resMatCol)
		// log.Printf("matRRStack: %.2g", mat.Formatted(matRRStack, mat.FormatMATLAB()))

		matCovSR := mat.NewSymDense(2, nil)
		stat.CovarianceMatrix(matCovSR, matSRStack.T(), nil)
		matCovRR := mat.NewSymDense(2, nil)
		stat.CovarianceMatrix(matCovRR, matRRStack.T(), nil)

		expCovSR := stat.Mean(matCovSR.RawSymmetric().Data, nil)
		expCovRR := stat.Mean(matCovRR.RawSymmetric().Data, nil)

		correlation = append(correlation, expCovSR/math.Sqrt(expCovRR*expCovSR))
	}

	log.Println(correlation)

	var (
		maxCorr float64
		maxCorrIdx int
	)
	maxCorr = math.Inf(-1)

	for i := 0; i < nodeCount; i++ {
		if correlation[i] > maxCorr {
			maxCorr = correlation[i]
			maxCorrIdx = i
		}
	}

	if maxCorr == math.Inf(-1) || maxCorr == math.Inf(1) {
		return -1, Error{"No leak found"}
	} 

	nodeLeaking := nodesId[maxCorrIdx]

	return nodeLeaking, nil
}