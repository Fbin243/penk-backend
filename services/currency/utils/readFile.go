package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type FishConfig struct {
	Type   string
	Number int
	Rate   float64
}

func LoadFishConfigs(filePath string) ([]FishConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	var configs []FishConfig
	reader := csv.NewReader(file)

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %v", err)
	}

	// Read records
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		number, _ := strconv.Atoi(record[1])
		rate, _ := strconv.ParseFloat(record[2], 64)
		configs = append(configs, FishConfig{
			Type:   record[0],
			Number: number,
			Rate:   rate,
		})
	}

	return configs, nil
}
