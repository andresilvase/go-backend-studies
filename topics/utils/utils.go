package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

type TableName string

const (
	Users    TableName = "users"
	Products TableName = "products"
)

func ProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)

		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}

		dir = parent
	}
}

func GetCSVFolderPath() (string, error) {
	root, err := ProjectRoot()

	if err != nil {
		return "", err
	}

	csvFolderPath := filepath.Join(
		root,
		"topics",
		"transactions",
		"data",
	)

	return csvFolderPath, nil
}

func CsvContent(csvFileName string) ([][]string, error) {
	csvFolderPath, err := GetCSVFolderPath()

	if err != nil {
		return [][]string{}, err
	}

	csvFile, err := os.Open(fmt.Sprintf("%s/%s.csv", csvFolderPath, csvFileName))

	if err != nil {
		return [][]string{}, err
	}

	defer csvFile.Close()

	fileReader := csv.NewReader(csvFile)

	csvRows, err := fileReader.ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return csvRows, nil
}
