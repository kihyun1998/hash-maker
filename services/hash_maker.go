package services

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kihyun1998/hash-maker/models"
)

const (
	SumFileName    = "hash_sum.txt"   // 해시 결과를 저장할 파일 이름
	ExecutableName = "hash-maker.exe" // 실행 파일 이름
)

// HashMaker는 해시 생성 작업을 담당하는 구조체
type HashMaker struct {
	config models.Config
	hashes map[string]models.FileInfo
}

// NewHashMaker는 새로운 HashMaker 인스턴스를 생성
func NewHashMaker(config models.Config) *HashMaker {
	return &HashMaker{
		config: config,
		hashes: make(map[string]models.FileInfo),
	}
}

// ProcessZipFile은 ZIP 파일에 대한 해시를 생성하고 추가
func (h *HashMaker) ProcessZipFile(zipFilePath string) error {
	// ZIP 파일 열기
	file, err := os.OpenFile(zipFilePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("ZIP 파일 열기 실패: %w", err)
	}
	defer file.Close()

	// 파일 정보 가져오기
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("파일 정보 가져오기 실패: %w", err)
	}

	// 원본 파일 크기 저장
	originalSize := fileInfo.Size()

	// 해시 계산
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Errorf("해시 계산 실패: %w", err)
	}

	// 해시값 생성
	hashSum := hash.Sum(nil)
	hashLength := len(hashSum)

	// 파일 끝으로 이동
	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		return fmt.Errorf("파일 포인터 이동 실패: %w", err)
	}

	// 해시 길이 쓰기 (4바이트)
	if err := binary.Write(file, binary.LittleEndian, int32(hashLength)); err != nil {
		return fmt.Errorf("해시 길이 쓰기 실패: %w", err)
	}

	// 해시값 쓰기
	if _, err := file.Write(hashSum); err != nil {
		return fmt.Errorf("해시값 쓰기 실패: %w", err)
	}

	fmt.Printf("해시 추가 완료\n")
	fmt.Printf("- 원본 크기: %d 바이트\n", originalSize)
	fmt.Printf("- 해시 길이: %d 바이트\n", hashLength)
	fmt.Printf("- 최종 크기: %d 바이트\n", originalSize+int64(hashLength)+4)

	return nil
}

// ProcessRequest는 설정에 따라 적절한 해시 처리를 수행
func (h *HashMaker) ProcessRequest() error {
	if h.config.UseZip {
		if h.config.ZipPath == "" {
			return fmt.Errorf("ZIP 모드에서는 ZIP 파일 경로가 필요합니다")
		}
		return h.ProcessZipFile(h.config.ZipPath)
	}

	if h.config.ZipFolder != "" && h.config.ZipName != "" {
		return h.CreateAndProcessZip()
	}

	return h.ProcessBasicHash()
}

// ProcessBasicHash는 기본 디렉토리 해싱을 처리
func (h *HashMaker) ProcessBasicHash() error {
	rootPath := h.config.StartPath
	if rootPath == "" {
		exePath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("실행 파일 경로 가져오기 실패: %w", err)
		}
		rootPath = filepath.Dir(exePath)
	}

	fmt.Printf("디렉토리 해시 계산 중: %s\n", rootPath)

	if err := h.calculateHashes(rootPath); err != nil {
		return fmt.Errorf("해시 계산 중 오류: %w", err)
	}

	sumFilePath := filepath.Join(rootPath, SumFileName)
	if err := h.writeSumFile(sumFilePath); err != nil {
		return fmt.Errorf("해시 파일 작성 중 오류: %w", err)
	}

	fmt.Printf("해시 파일 생성 완료: '%s'\n", sumFilePath)
	return nil
}
