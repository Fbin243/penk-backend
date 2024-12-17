package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// FishConfig represents configuration for each fish type.
type FishConfig struct {
	Type   string
	Number int
	Rate   float64
}

// ExchangeConfig represents configuration for exchanging fish for items.
type ExchangeConfig struct {
	ItemType string
	FishType string
	Number   int
	Increase int
}

// LoadFishConfigs loads fish configurations from a CSV file.
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

		number, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, fmt.Errorf("invalid number value in CSV: %v", err)
		}

		rate, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid rate value in CSV: %v", err)
		}

		configs = append(configs, FishConfig{
			Type:   record[0],
			Number: number,
			Rate:   rate,
		})
	}

	return configs, nil
}

// LoadExchangeConfigs loads exchange configurations from a CSV file.
func LoadExchangeConfigs(filePath string) ([]ExchangeConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %v", err)
	}

	var configs []ExchangeConfig
	// Read records
	for {
		record, err := reader.Read()
		if err != nil {
			break // End of file or error
		}

		number, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, fmt.Errorf("invalid number value in CSV: %v", err)
		}

		increase, err := strconv.Atoi(record[3])
		if err != nil {
			return nil, fmt.Errorf("invalid increase value in CSV: %v", err)
		}

		configs = append(configs, ExchangeConfig{
			ItemType: record[0],
			FishType: record[1],
			Number:   number,
			Increase: increase,
		})
	}

	return configs, nil
}
