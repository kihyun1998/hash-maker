package services

import (
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kihyun1998/hash-maker/models"
)

// calculatePathHash는 파일 경로의 해시값을 계산
func (h *HashMaker) calculatePathHash(relPath string) string {
	// 경로를 바이트로 변환하여 SHA-256 해시 계산
	hash := sha256.Sum256([]byte(relPath))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// calculateFileHash는 파일 내용의 해시값을 계산
func (h *HashMaker) calculateFileHash(filePath string) (string, error) {
	// 파일 열기
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("파일 열기 실패 (%s): %w", filePath, err)
	}
	defer file.Close()

	// SHA-256 해시 계산
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("파일 해시 계산 실패 (%s): %w", filePath, err)
	}

	// Base64 인코딩하여 반환
	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}

// writeSumFile은 계산된 해시값들을 파일에 작성
func (h *HashMaker) writeSumFile(sumFilePath string) error {
	// 파일 생성
	file, err := os.Create(sumFilePath)
	if err != nil {
		return fmt.Errorf("해시 파일 생성 실패: %w", err)
	}
	defer file.Close()

	// 버퍼 writer 생성
	writer := bufio.NewWriter(file)

	// 경로 정렬을 위한 슬라이스
	var paths []string
	for path := range h.hashes {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	// 정렬된 순서대로 해시 정보 쓰기
	for _, path := range paths {
		info := h.hashes[path]
		line := fmt.Sprintf("%s;%s;%s\n", info.FileType, info.PathHash, info.DataHash)

		if _, err := writer.WriteString(line); err != nil {
			return fmt.Errorf("해시 정보 쓰기 실패: %w", err)
		}
	}

	// 버퍼 플러시
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("파일 버퍼 플러시 실패: %w", err)
	}

	fmt.Printf("해시 파일 작성 완료: %s\n", sumFilePath)
	return nil
}

// calculateHashes는 디렉토리 내 모든 파일의 해시를 계산
func (h *HashMaker) calculateHashes(rootPath string) error {
	return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("파일 순회 중 오류: %w", err)
		}

		// 상대 경로 계산
		relPath, err := filepath.Rel(rootPath, path)
		if err != nil {
			return fmt.Errorf("상대 경로 계산 실패: %w", err)
		}

		// 경로 구분자를 슬래시로 통일
		relPath = filepath.ToSlash(relPath)

		// 해시 계산에서 제외할 파일들 필터링
		if relPath != SumFileName &&
			!strings.HasPrefix(relPath, ExecutableName) &&
			!info.IsDir() {

			// 경로 해시 계산
			pathHash := h.calculatePathHash(relPath)

			// 파일 내용 해시 계산
			dataHash, err := h.calculateFileHash(path)
			if err != nil {
				return fmt.Errorf("파일 해시 계산 실패 (%s): %w", path, err)
			}

			// 해시 정보 저장
			h.hashes[relPath] = models.FileInfo{
				FileType: "f",
				DataHash: dataHash,
				PathHash: pathHash,
			}

			fmt.Printf("파일 해시 계산 완료: %s\n", relPath)
		}

		return nil
	})
}
