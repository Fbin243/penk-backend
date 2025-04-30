package utils

import (
	"log"
	"os"
	"strings"
)

var root string

func init() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	root = strings.Split(dir, "tenk-backend")[0] + "tenk-backend"
}

func GetRoot() string {
	return root
}
