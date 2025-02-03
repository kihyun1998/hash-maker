package repository

import "github.com/kihyun1998/hash-maker/internal/domain/model"

// IArchiver는 압축 관련 작업을 위한 인터페이스
type IArchiver interface {
	CreateArchive(sourcePath string, targetPath string) error
	AddHashToArchive(archivePath string, hash []byte) error
	GetArchiveInfo(path string) (model.FileMetadata, error)
}
