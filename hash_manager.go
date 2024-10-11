package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type fileInfo struct {
	fileType string
	dataHash string
	pathHash string
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
			var dataHash, pathHash string
			var fileType string

			pathHash = calculatePathHash(relPath)

			if info.IsDir() {
				fileType = "d"
				dataHash, err = calculateDirHash(path)
			} else {
				fileType = "f"
				dataHash, err = calculateFileHash(path)
			}
			if err != nil {
				return err
			}
			hashes[relPath] = fileInfo{fileType: fileType, dataHash: dataHash, pathHash: pathHash}
		}

		return nil
	})

	return hashes, err
}

func calculatePathHash(relPath string) string {
	hash := sha256.New()
	hash.Write([]byte(relPath))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func calculateDirHash(dirPath string) (string, error) {
	hash := sha256.New()

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		subRelPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		// 디렉토리 자체와 제외할 파일들은 건너뜁니다
		if subRelPath == "." || subRelPath == sumFileName || subRelPath == executableName {
			return nil
		}

		// 파일 이름과 크기를 해시에 추가
		fmt.Fprintf(hash, "%s%d", subRelPath, info.Size())

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(hash, file); err != nil {
				return err
			}
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

	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
