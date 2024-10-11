package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type fileInfo struct {
	fileType string
	hash     string
}

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

func calculateHashes(rootPath string) (map[string]fileInfo, error) {
	hashes := make(map[string]fileInfo)

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(rootPath, path)
		if err != nil {
			return err
		}

		// 해시 계산에서 제외할 파일들
		if relPath != sumFileName && relPath != executableName {
			var hash string
			var fileType string

			if info.IsDir() {
				fileType = "d"
				hash, err = calculateDirHash(path)
			} else {
				fileType = "f"
				hash, err = calculateFileHash(path)
			}
			if err != nil {
				return err
			}
			hashes[relPath] = fileInfo{fileType: fileType, hash: hash}
		}

		return nil
	})

	return hashes, err
}

func calculateDirHash(dirPath string) (string, error) {
	hash := sha256.New()

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		// 디렉토리 자체와 제외할 파일들은 건너뜁니다
		if relPath == "." || relPath == sumFileName || relPath == executableName {
			return nil
		}

		// 파일 이름과 크기를 해시에 추가
		fmt.Fprintf(hash, "%s%d", relPath, info.Size())

		if !info.IsDir() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			hash.Write(data)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func writeSumFile(hashes map[string]fileInfo, sumFilePath string) error {
	file, err := os.Create(sumFilePath)
	if err != nil {
		return fmt.Errorf("failed to create sum file: %v", err)
	}
	defer file.Close()

	var paths []string
	for path := range hashes {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for _, path := range paths {
		info := hashes[path]
		_, err := fmt.Fprintf(file, "%s;%s;%s\n", info.fileType, strings.ReplaceAll(path, "\\", "/"), info.hash)
		if err != nil {
			return fmt.Errorf("failed to write to sum file: %v", err)
		}
	}

	return nil
}
