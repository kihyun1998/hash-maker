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
	// 0. 폴더 검사
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("출력 디렉토리 생성 실패: %w", err)
	}

	// 1. ZIP 파일 생성
	zipfile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("ZIP 파일 생성 실패: %w", err)
	}
	defer zipfile.Close()

	// 2. ZIP 작성자 생성
	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	// 3. 디렉토리 순회 시작
	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 4. ZIP 파일 헤더 생성
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("ZIP 헤더 생성 실패: %w", err)
		}

		// 5. 상대 경로 계산
		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return fmt.Errorf("상대 경로 계산 실패: %w", err)
		}

		// 6. 헤더 이름 설정 (경로 구분자를 슬래시로 통일)
		header.Name = filepath.ToSlash(relPath)
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate // 압축 방식 설정
		}

		// 7. ZIP 엔트리 생성
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("ZIP 엔트리 생성 실패: %w", err)
		}

		// 디렉토리면 더 이상 처리하지 않음
		if info.IsDir() {
			return nil
		}

		// 8. 파일 내용 복사
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("파일 열기 실패: %w", err)
		}
		defer file.Close()

		// 파일 내용을 ZIP에 복사
		_, err = io.Copy(writer, file)
		return err
	})
}
func (a *ZipArchiver) AddHashToArchive(archivePath string, hash []byte) error {
	// 1. ZIP 파일을 읽기/쓰기 모드로 열기 (끝에 추가 가능하도록)
	file, err := os.OpenFile(archivePath, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("ZIP 파일 열기 실패: %w", err)
	}
	defer file.Close()

	// 2. 해시 길이를 바이너리로 쓰기 (4바이트 정수)
	hashLength := int32(len(hash))
	if err := binary.Write(file, binary.LittleEndian, hashLength); err != nil {
		return fmt.Errorf("해시 길이 쓰기 실패: %w", err)
	}

	// 3. 실제 해시값 쓰기
	if _, err := file.Write(hash); err != nil {
		return fmt.Errorf("해시값 쓰기 실패: %w", err)
	}

	return nil
}

func (a *ZipArchiver) GetArchiveInfo(path string) (model.FileMetadata, error) {
	// 1. 파일 정보 가져오기
	info, err := os.Stat(path)
	if err != nil {
		return model.FileMetadata{}, fmt.Errorf("아카이브 정보 조회 실패: %w", err)
	}

	// 2. FileMetadata 구조체로 변환하여 반환
	return model.FileMetadata{
		RelativePath: filepath.Base(path), // 파일명만 추출
		Size:         info.Size(),         // 파일 크기
		IsDirectory:  false,               // ZIP 파일은 디렉토리가 아님
	}, nil
}
