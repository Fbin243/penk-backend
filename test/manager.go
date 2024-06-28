package test

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
)

type TestContext struct {
	IdCharacter    string `json:"idCharacter"`
	IdCustomMetric string `json:"idCustomMetric"`
	IdTimeTracking string `json:"idTimeTracking"`
}

type TestManager struct {
	fileName string
}

func NewTestManager() *TestManager {
	path, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %s", err)
	}

	fileName := strings.Replace(path, "/test", "", 1) + "/context.json"

	var file *os.File
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		file, err = os.Create(fileName)
		if err != nil {
			log.Fatalf("Failed to create file: %s", err)
		}
	}

	defer file.Close()

	return &TestManager{fileName: fileName}
}

func (m *TestManager) InitContext() {
	file, err := os.OpenFile(m.fileName, os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}

	defer file.Close()

	jsonCtx, err := json.Marshal(TestContext{
		IdCharacter:    "",
		IdCustomMetric: "",
		IdTimeTracking: "",
	})
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %s", err)
	}

	_, err = file.Write(jsonCtx)
	if err != nil {
		log.Fatalf("Failed to write file: %s", err)
	}
}

func (m *TestManager) GetContext() TestContext {
	file, err := os.Open(m.fileName)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}

	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	var ctx TestContext
	err = json.Unmarshal(byteValue, &ctx)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %s", err)
	}

	return ctx
}

func (m *TestManager) UpdateContext(newCtx TestContext) {
	file, err := os.OpenFile(m.fileName, os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}

	defer file.Close()

	jsonCtx, err := json.Marshal(newCtx)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %s", err)
	}

	_, err = file.Write(jsonCtx)
	if err != nil {
		log.Fatalf("Failed to write file: %s", err)
	}
}
