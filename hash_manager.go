package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

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
