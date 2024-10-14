package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

		relPath = filepath.ToSlash(relPath)

		if relPath != sumFileName && !strings.HasPrefix(relPath, executableName) && !info.IsDir() {
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
	fmt.Println(relPath)
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
