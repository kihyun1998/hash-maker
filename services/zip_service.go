package services

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CreateAndProcessZip은 디렉토리를 ZIP으로 압축하고 해시를 생성
func (h *HashMaker) CreateAndProcessZip() error {
	outputPath := h.config.ZipOutputPath
	if outputPath == "" {
		outputPath = "."
	}

	zipFilePath := filepath.Join(outputPath, h.config.ZipName)
	if err := h.zipDirectory(h.config.ZipFolder, zipFilePath); err != nil {
		return fmt.Errorf("ZIP 파일 생성 실패: %w", err)
	}

	fmt.Printf("ZIP 파일 생성 완료: %s\n", zipFilePath)
	return h.ProcessZipFile(zipFilePath)
}

// zipDirectory는 디렉토리를 ZIP 파일로 압축
func (h *HashMaker) zipDirectory(sourceDir, zipFile string) error {
	// ZIP 파일 생성을 위한 디렉토리 생성
	if err := os.MkdirAll(filepath.Dir(zipFile), os.ModePerm); err != nil {
		return fmt.Errorf("ZIP 디렉토리 생성 실패: %w", err)
	}

	// ZIP 파일 생성
	zipfile, err := os.Create(zipFile)
	if err != nil {
		return fmt.Errorf("ZIP 파일 생성 실패: %w", err)
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	// 디렉토리 순회하며 파일 압축
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// ZIP 엔트리 헤더 생성
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("ZIP 헤더 생성 실패: %w", err)
		}

		// 상대 경로 계산
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("상대 경로 계산 실패: %w", err)
		}

		header.Name = filepath.ToSlash(relPath)
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("ZIP 엔트리 생성 실패: %w", err)
		}

		if info.IsDir() {
			return nil
		}

		// 파일 내용 복사
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("파일 열기 실패: %w", err)
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})
}
