package repository

import "github.com/kihyun1998/hash-maker/internal/domain/model"

// IHashGenerator는 해시 생성을 위한 인터페이스
type IHashGenerator interface {
	GeneratePathHash(path string) (string, error)
	GenerateDataHash(data []byte) (string, error)
}

// IHashWriter는 해시 결과를 저장하기 위한 인터페이스
type IHashWriter interface {
	WriteHashSummary(hashes map[string]model.FileHash, outputPath string) error
}
