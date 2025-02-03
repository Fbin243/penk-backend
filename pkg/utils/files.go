package utils

import (
	"log"
	"os"
	"strings"
)

func GetRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return strings.Split(dir, "tenk-backend")[0] + "tenk-backend"
}
