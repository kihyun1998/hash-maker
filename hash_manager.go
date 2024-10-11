package main

import (
	"crypto/sha256"
	"encoding/base64"
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
		if relPath != sumFileName && relPath != executableName && !info.IsDir() {
			pathHash := calculatePathHash(relPath)

			dataHash, err := calculateFileHash(path)
			if err != nil {
				return err
			}
			hashes[relPath] = fileInfo{fileType: "f", dataHash: dataHash, pathHash: pathHash}

		}

		return nil
	})

	return hashes, err
}

func calculatePathHash(relPath string) string {
	hash := sha256.Sum256([]byte(relPath))
	return base64.StdEncoding.EncodeToString(hash[:])
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

	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}
