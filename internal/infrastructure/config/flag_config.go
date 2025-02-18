package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// FlagConfig는 명령줄 플래그 기반 설정 구현체
type FlagConfig struct {
	startPath     string
	zipPath       string
	useHash       bool
	zipFolder     string
	zipName       string
	zipOutputPath string
}

// NewFlagConfig는 새로운 FlagConfig 인스턴스를 생성하고 플래그를 파싱
func NewFlagConfig() (*FlagConfig, error) {
	config := &FlagConfig{}

	// 플래그 정의
	flag.StringVar(&config.startPath, "startPath", "", "시작 경로 지정")
	flag.StringVar(&config.zipPath, "zipPath", "", "ZIP 파일 경로 지정, -hash와 같이 사용한다.")
	flag.BoolVar(&config.useHash, "hash", false, "Hash 모드 사용, -zipPath로 지정한 zip파일의 해시값을 계산하여 포함한다.")
	flag.StringVar(&config.zipFolder, "zipfolder", "", "ZIP으로 만들 폴더 지정")
	flag.StringVar(&config.zipName, "zipname", "", "생성할 ZIP 파일 이름")
	flag.StringVar(&config.zipOutputPath, "zipoutput", "", "ZIP 파일 출력 경로")

	flag.Parse()

	// 기본값 설정
	if config.startPath == "" {
		exePath, err := os.Executable()
		if err != nil {
			return nil, err
		}
		config.startPath = filepath.Dir(exePath)
	}

	if config.zipOutputPath == "" {
		config.zipOutputPath = "."
	}

	return config, nil
}

// 인터페이스 구현
func (c *FlagConfig) GetStartPath() string  { return c.startPath }
func (c *FlagConfig) GetZipPath() string    { return c.zipPath }
func (c *FlagConfig) IsHashMode() bool      { return c.useHash }
func (c *FlagConfig) GetZipFolder() string  { return c.zipFolder }
func (c *FlagConfig) GetZipName() string    { return c.zipName }
func (c *FlagConfig) GetOutputPath() string { return c.zipOutputPath }

// Validate는 설정값의 유효성을 검증
func (c *FlagConfig) Validate() error {
	if c.useHash && c.zipPath == "" {
		return fmt.Errorf("ZIP 모드에서는 ZIP 파일 경로가 필요합니다")
	}
	return nil
}
