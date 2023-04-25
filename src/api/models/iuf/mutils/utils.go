package mutils

import (
	"errors"
	"io"
	"os"
	"strings"
)

const yamlFileDelimiter string = "---"

// Function to check a path exist or not
func IsPathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	} else {
		return true
	}
}

// Function to check a directory is empty or not
func IsEmptyDirectory(directoryPath string) bool {
	f, err := os.Open(directoryPath) // reading the path given

	if err != nil {
		return true
	}

	defer f.Close() // close the file handler post function exist

	_, err = f.Readdir(1) // checking for atleast single file in the directory

	if err == io.EOF { // if empty
		return true
	} else {
		return false
	}
}

// Function for string search operations
func StringFoundInArray(searchArray []string, searchString string) (found bool, index int) {
	found = false
	for i, x := range searchArray {
		if x == searchString {
			found = true
			index = i

			break
		}
	}
	return found, index
}

func Delete(orig []string, index int) ([]string, error) {
	if index < 0 || index >= len(orig) {
		return nil, errors.New("Index cannot be less than 0")
	}

	orig = append(orig[:index], orig[index+1:]...)

	return orig, nil
}

// File to manage
func ReadYamFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

// Function to split multi doc yaml
func SplitMultiYamlFile(fileData []byte) [][]byte {
	var yamlDataBytes [][]byte
	for _, yamlData := range strings.Split(string(fileData), yamlFileDelimiter) {
		if yamlData == "\n" || yamlData == "" { // skipping new line characters
			continue
		}
		yamlDataBytes = append(yamlDataBytes, []byte(yamlData))
	}

	return yamlDataBytes
}
