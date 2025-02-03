# hash-maker
## Project Structure

```
hash-maker/
├── internal/
    ├── domain/
    │   ├── model/
    │   │   └── file.go
    │   └── repository/
    │   │   ├── archiver.go
    │   │   ├── config.go
    │   │   ├── filesystem.go
    │   │   └── hasher.go
    ├── infrastructure/
    │   ├── fsys/
    │   │   └── local_filesystem.go
    │   ├── hashgen/
    │   │   └── sha256_generator.go
    │   └── ziputil/
    │   │   └── zip_archiver.go
    └── service/
    │   ├── archiver/
    │       └── archive_service.go
    │   └── hasher/
    │       └── hash_service.go
├── README.md
└── main.go
```

## README.md
```md
# hash-maker

파일 시스템의 무결성을 보장하기 위한 해시 생성 도구입니다. 일반 파일들의 해시를 생성하거나, ZIP 파일에 대한 해시를 추가할 수 있습니다.

## 주요 기능

- 디렉토리 내 파일들의 해시 생성
- ZIP 파일 생성 및 해시 추가
- 기존 ZIP 파일에 해시 추가

## 설치 방법

```bash
git clone https://github.com/kihyun1998/hash-maker.git
cd hash-maker
go build -o hash-maker
```

## 사용 방법

### 1. 디렉토리 해시 생성
```bash
./hash-maker -startPath "대상/디렉토리/경로"
```

Flutter 빌드 결과물의 경우:
```bash
./hash-maker -startPath "프로젝트명/build/windows/x64/runner/Release/."
```

### 2. ZIP 파일 해시 추가
```bash
./hash-maker -zip -zipPath "대상.zip"
```

### 3. 디렉토리를 ZIP으로 압축하고 해시 추가
```bash
./hash-maker -zipfolder "대상/디렉토리" -zipname "결과.zip" -zipoutput "출력/경로"
```

## 명령줄 옵션

| 옵션 | 설명 | 기본값 |
|------|------|--------|
| -startPath | 해시를 생성할 디렉토리 경로 | 실행 파일 위치 |
| -zip | ZIP 모드 사용 여부 | false |
| -zipPath | 처리할 ZIP 파일 경로 | - |
| -zipfolder | ZIP으로 만들 폴더 경로 | - |
| -zipname | 생성할 ZIP 파일 이름 | - |
| -zipoutput | ZIP 파일 출력 경로 | 현재 디렉토리 |

## 프로젝트 구조

```
hash-maker/
├── main.go                 # 애플리케이션 진입점
├── internal/
│   ├── domain/                 # 도메인 모델과 인터페이스
│   │   ├── model/             # 도메인 모델
│   │   └── repository/        # 저장소 인터페이스
│   ├── infrastructure/        # 인프라스트럭처 구현
│   │   ├── fsys/             # 파일 시스템 구현
│   │   ├── hashgen/          # 해시 생성 구현
│   │   └── ziputil/          # ZIP 파일 처리 구현
│   └── service/              # 비즈니스 로직
│       ├── hasher/           # 해시 처리 서비스
│       └── archiver/         # 압축 파일 처리 서비스
└── README.md
```

## 해시 파일 형식

생성되는 hash_sum.txt 파일은 다음과 같은 형식을 가집니다:
```
파일타입;경로해시;데이터해시
```

- 파일타입: 파일의 종류 (f: 일반 파일)
- 경로해시: 파일 경로의 SHA-256 해시 (Base64 인코딩)
- 데이터해시: 파일 내용의 SHA-256 해시 (Base64 인코딩)

## 기여하기

버그 리포트, 기능 요청, 풀 리퀘스트를 환영합니다.

## 라이선스

MIT License
```

이제 각 패키지와 주요 함수에 대한 주석을 추가하겠습니다. 예시로 몇 개만 보여드리겠습니다:

```go
// Package domain/model은 해시 메이커의 핵심 도메인 모델을 정의합니다.
// 이 패키지는 비즈니스 로직에서 사용되는 기본적인 데이터 구조와 타입을 포함합니다.
package model

// FileMetadata는 파일의 메타데이터를 나타내는 값 객체입니다.
// 파일 시스템에서 파일의 기본 정보를 저장하고 관리합니다.
type FileMetadata struct {
    // RelativePath는 기준 디렉토리로부터의 상대 경로입니다.
    RelativePath string

    // Size는 파일의 크기(바이트)입니다.
    Size int64

    // IsDirectory는 해당 항목이 디렉토리인지 여부를 나타냅니다.
    IsDirectory bool
}

```
## internal/domain/model/file.go
```go
package model

// FileMetadata는 파일의 메타데이터를 나타내는 값 객체
type FileMetadata struct {
	RelativePath string
	Size         int64
	IsDirectory  bool
}

// FileHash는 파일의 해시 정보를 나타내는 값 객체
type FileHash struct {
	PathHash string
	DataHash string
	FileType string
}

// HashResult는 해시 처리 결과를 나타내는 값 객체
type HashResult struct {
	OriginalSize int64
	HashLength   int32
	FinalSize    int64
}

```
## internal/domain/repository/archiver.go
```go
package repository

import "github.com/kihyun1998/hash-maker/internal/domain/model"

// IArchiver는 압축 관련 작업을 위한 인터페이스
type IArchiver interface {
	CreateArchive(sourcePath string, targetPath string) error
	AddHashToArchive(archivePath string, hash []byte) error
	GetArchiveInfo(path string) (model.FileMetadata, error)
}

```
## internal/domain/repository/config.go
```go
package repository

// IConfigProvider는 설정 관리를 위한 인터페이스
type IConfigProvider interface {
	GetStartPath() string
	GetZipPath() string
	IsZipMode() bool
	GetZipFolder() string
	GetZipName() string
	GetOutputPath() string
}

```
## internal/domain/repository/filesystem.go
```go
package repository

import "github.com/kihyun1998/hash-maker/internal/domain/model"

// IFileSystem은 파일 시스템 작업을 위한 인터페이스
type IFileSystem interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
	WalkDirectory(root string, callback func(model.FileMetadata) error) error
	GetFileInfo(path string) (model.FileMetadata, error)
}

```
## internal/domain/repository/hasher.go
```go
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

```
## internal/infrastructure/fsys/local_filesystem.go
```go
package fsys

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kihyun1998/hash-maker/internal/domain/model"
	"github.com/kihyun1998/hash-maker/internal/domain/repository"
)

// LocalFileSystem은 로컬 파일 시스템 구현체
type LocalFileSystem struct{}

// NewLocalFileSystem은 새로운 LocalFileSystem 인스턴스를 생성
func NewLocalFileSystem() repository.IFileSystem {
	return &LocalFileSystem{}
}

func (fs *LocalFileSystem) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (fs *LocalFileSystem) WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func (fs *LocalFileSystem) WalkDirectory(root string, callback func(model.FileMetadata) error) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("디렉토리 순회 중 오류: %w", err)
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("상대 경로 계산 실패: %w", err)
		}

		metadata := model.FileMetadata{
			RelativePath: filepath.ToSlash(relPath),
			Size:         info.Size(),
			IsDirectory:  info.IsDir(),
		}

		return callback(metadata)
	})
}

func (fs *LocalFileSystem) GetFileInfo(path string) (model.FileMetadata, error) {
	info, err := os.Stat(path)
	if err != nil {
		return model.FileMetadata{}, fmt.Errorf("파일 정보 조회 실패: %w", err)
	}

	return model.FileMetadata{
		RelativePath: filepath.Base(path),
		Size:         info.Size(),
		IsDirectory:  info.IsDir(),
	}, nil
}

```
## internal/infrastructure/hashgen/sha256_generator.go
```go
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

```
## internal/infrastructure/ziputil/zip_archiver.go
```go
package ziputil

import (
	"archive/zip"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kihyun1998/hash-maker/internal/domain/model"
	"github.com/kihyun1998/hash-maker/internal/domain/repository"
)

// ZipArchiver는 ZIP 압축 구현체
type ZipArchiver struct{}

// NewZipArchiver는 새로운 ZipArchiver 인스턴스를 생성
func NewZipArchiver() repository.IArchiver {
	return &ZipArchiver{}
}

func (a *ZipArchiver) CreateArchive(sourcePath string, targetPath string) error {
	zipfile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("ZIP 파일 생성 실패: %w", err)
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("ZIP 헤더 생성 실패: %w", err)
		}

		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return fmt.Errorf("상대 경로 계산 실패: %w", err)
		}

		header.Name = filepath.ToSlash(relPath)
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("ZIP 엔트리 생성 실패: %w", err)
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("파일 열기 실패: %w", err)
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})
}

func (a *ZipArchiver) AddHashToArchive(archivePath string, hash []byte) error {
	file, err := os.OpenFile(archivePath, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("ZIP 파일 열기 실패: %w", err)
	}
	defer file.Close()

	hashLength := int32(len(hash))
	if err := binary.Write(file, binary.LittleEndian, hashLength); err != nil {
		return fmt.Errorf("해시 길이 쓰기 실패: %w", err)
	}

	if _, err := file.Write(hash); err != nil {
		return fmt.Errorf("해시값 쓰기 실패: %w", err)
	}

	return nil
}

func (a *ZipArchiver) GetArchiveInfo(path string) (model.FileMetadata, error) {
	info, err := os.Stat(path)
	if err != nil {
		return model.FileMetadata{}, fmt.Errorf("아카이브 정보 조회 실패: %w", err)
	}

	return model.FileMetadata{
		RelativePath: filepath.Base(path),
		Size:         info.Size(),
		IsDirectory:  false,
	}, nil
}

```
## internal/service/archiver/archive_service.go
```go
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

```
## internal/service/hasher/hash_service.go
```go
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

```
## main.go
```go
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kihyun1998/hash-maker/internal/domain/model"
	"github.com/kihyun1998/hash-maker/internal/infrastructure/fsys"
	"github.com/kihyun1998/hash-maker/internal/infrastructure/hashgen"
	"github.com/kihyun1998/hash-maker/internal/infrastructure/ziputil"
	"github.com/kihyun1998/hash-maker/internal/service/archiver"
	"github.com/kihyun1998/hash-maker/internal/service/hasher"
)

// Config는 명령줄 인자를 관리하는 구조체
type Config struct {
	StartPath     string
	ZipPath       string
	UseZip        bool
	ZipFolder     string
	ZipName       string
	ZipOutputPath string
}

// ConfigProvider는 Config 구조체를 기반으로 한 설정 제공자
type ConfigProvider struct {
	config Config
}

func (p *ConfigProvider) GetStartPath() string  { return p.config.StartPath }
func (p *ConfigProvider) GetZipPath() string    { return p.config.ZipPath }
func (p *ConfigProvider) IsZipMode() bool       { return p.config.UseZip }
func (p *ConfigProvider) GetZipFolder() string  { return p.config.ZipFolder }
func (p *ConfigProvider) GetZipName() string    { return p.config.ZipName }
func (p *ConfigProvider) GetOutputPath() string { return p.config.ZipOutputPath }

func main() {
	// 설정 파싱
	config := parseFlags()
	configProvider := &ConfigProvider{config: config}

	// 의존성 초기화
	hashGenerator := hashgen.NewSHA256Generator()
	fileSystem := fsys.NewLocalFileSystem()
	zipArchiver := ziputil.NewZipArchiver()

	// 서비스 초기화
	hashService := hasher.NewHashService(hashGenerator, fileSystem, configProvider)
	archiveService := archiver.NewArchiveService(zipArchiver, hashGenerator, fileSystem, configProvider)

	// 요청 처리
	if err := processRequest(config, hashService, archiveService); err != nil {
		fmt.Fprintf(os.Stderr, "오류 발생: %v\n", err)
		os.Exit(1)
	}
}

// parseFlags는 명령줄 인자를 파싱하여 Config 구조체를 반환
func parseFlags() Config {
	config := Config{}

	flag.StringVar(&config.StartPath, "startPath", "", "시작 경로 지정")
	flag.StringVar(&config.ZipPath, "zipPath", "", "ZIP 파일 경로 지정")
	flag.BoolVar(&config.UseZip, "zip", false, "ZIP 모드 사용")
	flag.StringVar(&config.ZipFolder, "zipfolder", "", "ZIP으로 만들 폴더 지정")
	flag.StringVar(&config.ZipName, "zipname", "", "생성할 ZIP 파일 이름")
	flag.StringVar(&config.ZipOutputPath, "zipoutput", "", "ZIP 파일 출력 경로")

	flag.Parse()

	return config
}

// processRequest는 설정에 따라 적절한 서비스를 호출
func processRequest(
	config Config,
	hashService *hasher.HashService,
	archiveService *archiver.ArchiveService,
) error {
	// ZIP 모드일 경우
	if config.UseZip {
		if config.ZipPath == "" {
			return fmt.Errorf("ZIP 모드에서는 ZIP 파일 경로가 필요합니다")
		}
		result, err := archiveService.ProcessArchive(config.ZipPath)
		if err != nil {
			return err
		}
		printHashResult(result)
		return nil
	}

	// ZIP 생성 모드일 경우
	if config.ZipFolder != "" && config.ZipName != "" {
		outputPath := config.ZipOutputPath
		if outputPath == "" {
			outputPath = "."
		}
		zipPath := filepath.Join(outputPath, config.ZipName)
		result, err := archiveService.CreateAndProcessArchive(config.ZipFolder, zipPath)
		if err != nil {
			return err
		}
		printHashResult(result)
		return nil
	}

	// 기본 해시 모드일 경우
	startPath := config.StartPath
	if startPath == "" {
		exePath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("실행 파일 경로 가져오기 실패: %w", err)
		}
		startPath = filepath.Dir(exePath)
	}

	return hashService.ProcessDirectory(startPath)
}

// printHashResult는 해시 처리 결과를 출력
func printHashResult(result *model.HashResult) {
	fmt.Printf("해시 추가 완료\n")
	fmt.Printf("- 원본 크기: %d 바이트\n", result.OriginalSize)
	fmt.Printf("- 해시 길이: %d 바이트\n", result.HashLength)
	fmt.Printf("- 최종 크기: %d 바이트\n", result.FinalSize)
}

```
