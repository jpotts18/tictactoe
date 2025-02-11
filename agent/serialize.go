package agent

import (
	"encoding/json"
	"os"
)

type QTableData struct {
	QTable map[string][]float64 `json:"qtable"`
}

type MonteCarloData struct {
	QTable  map[string][]float64         `json:"qtable"`
	Returns map[string]map[int][]float64 `json:"returns"`
}

func SaveQTable(filename string, data map[string][]float64) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(QTableData{QTable: data})
}

func LoadQTable(filename string) (map[string][]float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data QTableData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}
	return data.QTable, nil
} 
