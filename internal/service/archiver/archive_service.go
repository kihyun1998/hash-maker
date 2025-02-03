package archiver

import (
	"crypto/sha256"
	"fmt"

	"github.com/kihyun1998/hash-maker/internal/domain/model"
	"github.com/kihyun1998/hash-maker/internal/domain/repository"
)

// ArchiveService는 압축 파일 처리를 담당하는 서비스
type ArchiveService struct {
	archiver   repository.IArchiver
	hashGen    repository.IHashGenerator
	fileSystem repository.IFileSystem
	config     repository.IConfigProvider
}

// NewArchiveService는 새로운 ArchiveService 인스턴스를 생성
func NewArchiveService(
	archiver repository.IArchiver,
	hashGen repository.IHashGenerator,
	fs repository.IFileSystem,
	config repository.IConfigProvider,
) *ArchiveService {
	return &ArchiveService{
		archiver:   archiver,
		hashGen:    hashGen,
		fileSystem: fs,
		config:     config,
	}
}

// ProcessArchive는 압축 파일에 대한 해시 처리를 수행
func (s *ArchiveService) ProcessArchive(archivePath string) (*model.HashResult, error) {
	// 압축 파일 정보 조회
	fileInfo, err := s.archiver.GetArchiveInfo(archivePath)
	if err != nil {
		return nil, fmt.Errorf("압축 파일 정보 조회 실패: %w", err)
	}

	// 파일 데이터 읽기
	data, err := s.fileSystem.ReadFile(archivePath)
	if err != nil {
		return nil, fmt.Errorf("압축 파일 읽기 실패: %w", err)
	}

	// 해시 생성
	hash := sha256.Sum256(data)
	hashLength := int32(len(hash))

	// 해시 추가
	if err := s.archiver.AddHashToArchive(archivePath, hash[:]); err != nil {
		return nil, fmt.Errorf("해시 추가 실패: %w", err)
	}

	// 최종 파일 정보 조회
	finalInfo, err := s.archiver.GetArchiveInfo(archivePath)
	if err != nil {
		return nil, fmt.Errorf("최종 파일 정보 조회 실패: %w", err)
	}

	// 처리 결과 반환
	return &model.HashResult{
		OriginalSize: fileInfo.Size,
		HashLength:   hashLength,
		FinalSize:    finalInfo.Size,
	}, nil
}

// CreateAndProcessArchive는 디렉토리를 압축하고 해시를 추가
func (s *ArchiveService) CreateAndProcessArchive(sourcePath, targetPath string) (*model.HashResult, error) {
	// 압축 파일 생성
	if err := s.archiver.CreateArchive(sourcePath, targetPath); err != nil {
		return nil, fmt.Errorf("압축 파일 생성 실패: %w", err)
	}

	// 생성된 압축 파일 처리
	return s.ProcessArchive(targetPath)
}
