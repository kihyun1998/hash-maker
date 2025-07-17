package platform

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/kihyun1998/hash-maker/internal/domain/repository"
)

// ExecutableProvider는 실행 파일 관련 정보를 제공하는 구현체
type ExecutableProvider struct{}

// NewExecutableProvider는 새로운 ExecutableProvider 인스턴스를 생성
func NewExecutableProvider() repository.IPlatformProvider {
	return &ExecutableProvider{}
}

// GetExecutableName은 현재 플랫폼에 맞는 실행 파일명을 반환
func (p *ExecutableProvider) GetExecutableName() string {
	baseName := "hash-maker"
	if runtime.GOOS == "windows" {
		return baseName + ".exe"
	}
	return baseName
}

// GetCurrentExecutablePath는 현재 실행 중인 파일의 경로를 반환
func (p *ExecutableProvider) GetCurrentExecutablePath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exePath), nil
}
