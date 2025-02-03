package repository

import "github.com/kihyun1998/hash-maker/internal/domain/model"

// IFileSystem은 파일 시스템 작업을 위한 인터페이스
type IFileSystem interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
	WalkDirectory(root string, callback func(model.FileMetadata) error) error
	GetFileInfo(path string) (model.FileMetadata, error)
}
