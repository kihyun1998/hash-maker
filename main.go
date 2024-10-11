package main

import (
	"crypto/sha256"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var startPath string
var rootPath string
var zipPath string
var useZip bool

func init() {
	flag.StringVar(&startPath, "startPath", "", "The valiable startPath is start root path")
	flag.StringVar(&zipPath, "zipPath", "", "The path to the zip file to be hashed")
	flag.BoolVar(&useZip, "zip", false, "Use zip file mode")
	flag.Parse()
}

const (
	sumFileName    = "hash_sum.txt"
	executableName = "hash-maker.exe"
)

func main() {
	if useZip {
		if zipPath == "" {
			fmt.Println("Error: -zipPath is required when -zip flag is used")
			return
		}

		err := processZipFile(zipPath)
		if err != nil {
			fmt.Printf("Error processing zip file: %v\n", err)
			return
		}
		fmt.Println("Zip file processed successfully.")
		return
	}

	processBasicHash()
}

func processZipFile(zipFilePath string) error {
	file, err := os.OpenFile(zipFilePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	originalSize := fileInfo.Size()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Errorf("failed to calculate hash: %v", err)
	}

	hashSum := hash.Sum(nil)
	hashLength := len(hashSum)

	// 파일 포인터를 파일의 끝으로 이동합니다.
	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		return fmt.Errorf("failed to seek to end of file: %v", err)
	}

	// 해시 길이를 먼저 씁니다 (4바이트 정수).
	if err := binary.Write(file, binary.LittleEndian, int32(hashLength)); err != nil {
		return fmt.Errorf("failed to write hash length: %v", err)
	}

	// 해시를 파일의 끝에 씁니다.
	if _, err := file.Write(hashSum); err != nil {
		return fmt.Errorf("failed to write hash to file: %v", err)
	}

	fmt.Printf("Hash appended to file. Original size: %d bytes, Hash length: %d bytes, New size: %d bytes\n",
		originalSize, hashLength, originalSize+int64(hashLength)+4)

	return nil
}

func processBasicHash() {
	if startPath != "" {
		rootPath = filepath.Dir(startPath)
	} else {
		exePath, err := os.Executable()
		if err != nil {
			fmt.Printf("Error getting executable path: %v\n", err)
			return
		}

		rootPath = filepath.Dir(exePath)
	}

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
