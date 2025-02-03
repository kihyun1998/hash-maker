package ziputil

import (
	"archive/zip"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kihyun1998/hash-maker/internal/domain/model"
	"github.com/kihyun1998/hash-maker/internal/domain/repository"
)

// ZipArchiver는 ZIP 압축 구현체
type ZipArchiver struct{}

// NewZipArchiver는 새로운 ZipArchiver 인스턴스를 생성
func NewZipArchiver() repository.IArchiver {
	return &ZipArchiver{}
}

func (a *ZipArchiver) CreateArchive(sourcePath string, targetPath string) error {
	zipfile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("ZIP 파일 생성 실패: %w", err)
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("ZIP 헤더 생성 실패: %w", err)
		}

		relPath, err := filepath.Rel(sourcePath, path)
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

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("파일 열기 실패: %w", err)
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})
}

func (a *ZipArchiver) AddHashToArchive(archivePath string, hash []byte) error {
	file, err := os.OpenFile(archivePath, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("ZIP 파일 열기 실패: %w", err)
	}
	defer file.Close()

	hashLength := int32(len(hash))
	if err := binary.Write(file, binary.LittleEndian, hashLength); err != nil {
		return fmt.Errorf("해시 길이 쓰기 실패: %w", err)
	}

	if _, err := file.Write(hash); err != nil {
		return fmt.Errorf("해시값 쓰기 실패: %w", err)
	}

	return nil
}

func (a *ZipArchiver) GetArchiveInfo(path string) (model.FileMetadata, error) {
	info, err := os.Stat(path)
	if err != nil {
		return model.FileMetadata{}, fmt.Errorf("아카이브 정보 조회 실패: %w", err)
	}

	return model.FileMetadata{
		RelativePath: filepath.Base(path),
		Size:         info.Size(),
		IsDirectory:  false,
	}, nil
}
