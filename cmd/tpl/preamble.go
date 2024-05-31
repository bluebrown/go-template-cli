package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// only this number of bytes in the header is scanned for preamble
// before loading the whole file
const PreambleBufferSize = 512
const PreambleSymbol = "-|-"
const PreambleMagicString = "-|- build-edge:"

func doesContainPreamble(outputPath string) (bool, error) {
	file, err := os.Open(outputPath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	buffer := make([]byte, PreambleBufferSize)
	n, err := file.Read(buffer)
	if err != nil {
		return false, err
	}

	if n == 0 {
		return false, nil
	}

	fileHeader := string(buffer)
	return strings.Contains(fileHeader, PreambleMagicString), nil
}

func GetPreamble(outputPath string) (string, error) {
	containsPreamble, err := doesContainPreamble(outputPath)
	if err != nil {
		return "", err
	}

	if !containsPreamble {
		return "", nil
	}

	file, err := os.Open(outputPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var outBuf string

	// scan until first build edge spec line
	found := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, PreambleSymbol) {
			outBuf += line + "\n"
			found = true
			break
		}

		outBuf += line + "\n"
	}

	if !found {
		return "", fmt.Errorf("Build spec scan error")
	}

	// scan until last build edge spec line
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, PreambleSymbol) {
			break
		}

		outBuf += line + "\n"
	}

	return outBuf, nil
}
