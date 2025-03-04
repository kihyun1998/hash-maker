package hasher

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/kihyun1998/hash-maker/internal/domain/model"
	"github.com/kihyun1998/hash-maker/internal/domain/repository"
)

const (
	SumFileName    = "hash_sum.txt"   // 해시 결과를 저장할 파일 이름
	ExecutableName = "hash-maker.exe" // 실행 파일 이름
)

// HashService는 파일 해시 처리를 담당하는 서비스
type HashService struct {
	hashGenerator repository.IHashGenerator
	fileSystem    repository.IFileSystem
	config        repository.IConfigProvider
	hashes        map[string]model.FileHash
}

// NewHashService는 새로운 HashService 인스턴스를 생성
func NewHashService(
	hashGen repository.IHashGenerator,
	fs repository.IFileSystem,
	config repository.IConfigProvider,
) *HashService {
	return &HashService{
		hashGenerator: hashGen,
		fileSystem:    fs,
		config:        config,
		hashes:        make(map[string]model.FileHash),
	}
}

// ProcessDirectory는 디렉토리의 모든 파일에 대한 해시를 생성
func (s *HashService) ProcessDirectory(rootPath string) error {
	// 디렉토리 순회하며 해시 생성
	err := s.fileSystem.WalkDirectory(rootPath, func(metadata model.FileMetadata) error {
		// 해시 생성 제외 조건 체크
		if s.shouldSkipFile(metadata) {
			return nil
		}

		// 경로 해시 생성
		pathHash, err := s.hashGenerator.GeneratePathHash(metadata.RelativePath)
		if err != nil {
			return fmt.Errorf("경로 해시 생성 실패 (%s): %w", metadata.RelativePath, err)
		}

		// 파일 데이터 읽기
		fullPath := filepath.Join(rootPath, metadata.RelativePath)
		data, err := s.fileSystem.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("파일 읽기 실패 (%s): %w", fullPath, err)
		}

		// 데이터 해시 생성
		dataHash, err := s.hashGenerator.GenerateDataHash(data)
		if err != nil {
			return fmt.Errorf("데이터 해시 생성 실패 (%s): %w", fullPath, err)
		}

		// 해시 정보 저장
		s.hashes[metadata.RelativePath] = model.FileHash{
			FileType: "f",
			PathHash: pathHash,
			DataHash: dataHash,
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("디렉토리 처리 중 오류: %w", err)
	}

	// 해시 결과 파일 생성
	return s.writeHashSummary(rootPath)
}

// shouldSkipFile은 해시 생성을 건너뛸 파일인지 확인
func (s *HashService) shouldSkipFile(metadata model.FileMetadata) bool {
	return metadata.RelativePath == SumFileName ||
		metadata.RelativePath == ExecutableName ||
		metadata.IsDirectory
}

// writeHashSummary는 생성된 해시를 파일에 저장
func (s *HashService) writeHashSummary(rootPath string) error {
	// 결과를 저장할 파일 경로
	sumFilePath := filepath.Join(rootPath, SumFileName)

	// 정렬된 결과 생성
	var lines []string
	for _, hash := range s.hashes {
		line := fmt.Sprintf("%s;%s;%s", hash.FileType, hash.PathHash, hash.DataHash)
		lines = append(lines, line)
	}
	sort.Strings(lines)

	// 결과 저장
	content := []byte(formatHashSummary(lines))
	if err := s.fileSystem.WriteFile(sumFilePath, content); err != nil {
		return fmt.Errorf("해시 결과 파일 저장 실패: %w", err)
	}

	return nil
}

// formatHashSummary는 해시 결과를 포맷팅
func formatHashSummary(lines []string) string {
	result := ""
	for _, line := range lines {
		result += line + "\n"
	}
	return result
}
