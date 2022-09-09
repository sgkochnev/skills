package linksSearcher

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

const (
	emptyByte             = "\x00"
	maxResponseBufferSize = 128 * 1024
	permissions           = 666
)

func createInputFile(fileName, input string) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY, permissions)
	if os.IsNotExist(err) {
		file, err = os.Create(fileName)
	}
	if err != nil {
		return nil, fmt.Errorf("coulndn't create input file: %v", err)
	}
	_, err = file.WriteString(input)

	if err != nil {
		return nil, fmt.Errorf("coulndn't write to input file: %v", err)
	}
	return file, nil
}

func readOutputFile(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", fmt.Errorf("coulndn't open output file: %v", err)
	}

	var responseBytes []byte
	for {
		responseBuffer := make([]byte, maxResponseBufferSize)
		_, err = file.Read(responseBuffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error reading response: %v", err)
		}
		responseBytes = append(responseBytes, responseBuffer...)
	}
	responseBytes = bytes.Trim(responseBytes, emptyByte)

	return string(responseBytes), nil
}

func readOutputFile2(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return []string{}, fmt.Errorf("coulndn't open output file: %v", err)
	}
	reader := bufio.NewReader(file)
	var result []string
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return []string{}, fmt.Errorf("error reading response: %v", err)
		}
		result = append(result, string(line))
	}

	return result, nil
}

func deleteFiles(names ...string) error {
	var err error
	for i, name := range names {
		err = os.Remove(name)
		if err != nil {
			return fmt.Errorf("Error [%d %s: %w\n]", i, name, err)
		}
	}
	return err
}
