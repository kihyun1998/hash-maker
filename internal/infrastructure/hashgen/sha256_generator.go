package hashgen

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/kihyun1998/hash-maker/internal/domain/repository"
)

// SHA256Generator는 SHA-256 해시 생성기 구현체
type SHA256Generator struct{}

// NewSHA256Generator는 새로운 SHA256Generator 인스턴스를 생성
func NewSHA256Generator() repository.IHashGenerator {
	return &SHA256Generator{}
}

func (g *SHA256Generator) GeneratePathHash(path string) (string, error) {
	hash := sha256.Sum256([]byte(path))
	return base64.StdEncoding.EncodeToString(hash[:]), nil
}

func (g *SHA256Generator) GenerateDataHash(data []byte) (string, error) {
	hash := sha256.New()
	if _, err := hash.Write(data); err != nil {
		return "", fmt.Errorf("데이터 해시 생성 실패: %w", err)
	}
	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}
