package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	sumFileName    = "hash_sum.txt"
	executableName = "hash-maker.exe"
)

func main() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error getting executable path: %v\n", err)
		return
	}

	rootPath := filepath.Dir(exePath)

	fmt.Printf("Calculating hashes for directory: %s\n", rootPath)

	hashes, err := calculateHashes(rootPath)
	if err != nil {
		fmt.Printf("Error calculating hashes: %v\n", err)
		return
	}

	sumFilePath := filepath.Join(rootPath, sumFileName)
	err = writeSumFile(hashes, sumFilePath)
	if err != nil {
		fmt.Printf("Error writing sum file: %v\n", err)
		return
	}

	fmt.Printf("Hash sum file '%s' has been created successfully.\n", sumFilePath)
}
