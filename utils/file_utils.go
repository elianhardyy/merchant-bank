package utils

import (
	"encoding/json"
	"go-json/constant"
	"log"
	"os"
	"path/filepath"
)

func EnsureJSONFiles() {
	files := []string{constant.USER_FILE, constant.MERCHANT_FILE, constant.TRANSACTION_FILE}

	for _, filepath := range files {
		if !fileExists(filepath) {
			err := createEmptyJSONFile(filepath)
			if err != nil {
				log.Fatalf("Failed to create %s: %v", filepath, err)
			}
		}
	}
}

func fileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

func createEmptyJSONFile(filePath string) error {
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	var initialData interface{}
	switch filePath {
	case constant.USER_FILE:
		initialData = []map[string]interface{}{}
	case constant.MERCHANT_FILE:
		initialData = []map[string]interface{}{}
	case constant.TRANSACTION_FILE:
		initialData = []map[string]interface{}{}
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(initialData)
}

func ReadJSONFile(filepath string, v interface{}) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(v)
}

func WriteJSONFile(filepath string, v interface{}) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}
