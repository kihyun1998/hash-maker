# hash-maker
## Project Structure

```
hash-maker/
├── models/
    ├── config.go
    └── file_info.go
├── services/
    ├── hash_maker.go
    ├── hash_service.go
    └── zip_service.go
├── README.md
└── main.go
```

## README.md
```md
# hash-maker
 This is making hash windows app write by golang


## 사용방법 예시

### 일반 해시

```bash
.\hash-maker.exe -startPath C:\Users\User\study-flutter\update_test_app_1\build\windows\x64\runner\Release\.
```

flutter로 build windows했다면

`프로젝트명\build\windows\x64\runner\Release`

아래에 빌드 된다. 그렇기에

`프로젝트명\build\windows\x64\runner\Release\.` 이렇게 `\.` 이걸 추가해줘서 Release 안에를 hash해야한다.


### zip 해시

```bash
.\hash-maker.exe -zip -zipPath C:\Users\User\study-flutter\update_test_app_1\build\windows\x64\runner\Release.zip
```

이런식으로 zip파일 지정해서 해시를 뜰 수 있습니다.


## 사용 예시

```bash
# 기본 해시 생성
./hash-maker -startPath "C:\Project\build\Release"

# 특정 폴더를 ZIP으로 만들고 해시 생성
./hash-maker -zipfolder "C:\Project\build\Release" -zipname "update.zip"

# 기존 ZIP 파일에 해시 추가
./hash-maker -zip -zipPath "C:\Project\update.zip"
```
```
## main.go
```go
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kihyun1998/hash-maker/models"
	"github.com/kihyun1998/hash-maker/services"
)

func main() {
	// 명령줄 인자 파싱
	config := parseFlags()

	// HashMaker 인스턴스 생성
	hashMaker := services.NewHashMaker(config)

	// 해시 처리 실행
	if err := hashMaker.ProcessRequest(); err != nil {
		fmt.Fprintf(os.Stderr, "오류 발생: %v\n", err)
		os.Exit(1)
	}
}

// parseFlags는 명령줄 인자를 파싱하여 Config 구조체를 반환
func parseFlags() models.Config {
	config := models.Config{}

	flag.StringVar(&config.StartPath, "startPath", "", "시작 경로 지정")
	flag.StringVar(&config.ZipPath, "zipPath", "", "ZIP 파일 경로 지정")
	flag.BoolVar(&config.UseZip, "zip", false, "ZIP 모드 사용")
	flag.StringVar(&config.ZipFolder, "zipfolder", "", "ZIP으로 만들 폴더 지정")
	flag.StringVar(&config.ZipName, "zipname", "", "생성할 ZIP 파일 이름")
	flag.StringVar(&config.ZipOutputPath, "zipoutput", "", "ZIP 파일 출력 경로")

	flag.Parse()

	return config
}

```
## models/config.go
```go
package models

// Config는 명령줄 인자값들을 저장하는 구조체
type Config struct {
	StartPath     string // 시작 경로
	ZipPath       string // ZIP 파일 경로
	UseZip        bool   // ZIP 모드 사용 여부
	ZipFolder     string // ZIP으로 만들 폴더 경로
	ZipName       string // 생성할 ZIP 파일 이름
	ZipOutputPath string // ZIP 파일 출력 경로
}

```
## models/file_info.go
```go
package models

// FileInfo는 파일의 해시 정보를 저장하는 구조체
type FileInfo struct {
	FileType string // 파일 타입
	DataHash string // 파일 내용의 해시값
	PathHash string // 파일 경로의 해시값
}

```
## services/hash_maker.go
```go
package services

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kihyun1998/hash-maker/models"
)

const (
	SumFileName    = "hash_sum.txt"   // 해시 결과를 저장할 파일 이름
	ExecutableName = "hash-maker.exe" // 실행 파일 이름
)

// HashMaker는 해시 생성 작업을 담당하는 구조체
type HashMaker struct {
	config models.Config
	hashes map[string]models.FileInfo
}

// NewHashMaker는 새로운 HashMaker 인스턴스를 생성
func NewHashMaker(config models.Config) *HashMaker {
	return &HashMaker{
		config: config,
		hashes: make(map[string]models.FileInfo),
	}
}

// ProcessZipFile은 ZIP 파일에 대한 해시를 생성하고 추가
func (h *HashMaker) ProcessZipFile(zipFilePath string) error {
	// ZIP 파일 열기
	file, err := os.OpenFile(zipFilePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("ZIP 파일 열기 실패: %w", err)
	}
	defer file.Close()

	// 파일 정보 가져오기
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("파일 정보 가져오기 실패: %w", err)
	}

	// 원본 파일 크기 저장
	originalSize := fileInfo.Size()

	// 해시 계산
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Errorf("해시 계산 실패: %w", err)
	}

	// 해시값 생성
	hashSum := hash.Sum(nil)
	hashLength := len(hashSum)

	// 파일 끝으로 이동
	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		return fmt.Errorf("파일 포인터 이동 실패: %w", err)
	}

	// 해시 길이 쓰기 (4바이트)
	if err := binary.Write(file, binary.LittleEndian, int32(hashLength)); err != nil {
		return fmt.Errorf("해시 길이 쓰기 실패: %w", err)
	}

	// 해시값 쓰기
	if _, err := file.Write(hashSum); err != nil {
		return fmt.Errorf("해시값 쓰기 실패: %w", err)
	}

	fmt.Printf("해시 추가 완료\n")
	fmt.Printf("- 원본 크기: %d 바이트\n", originalSize)
	fmt.Printf("- 해시 길이: %d 바이트\n", hashLength)
	fmt.Printf("- 최종 크기: %d 바이트\n", originalSize+int64(hashLength)+4)

	return nil
}

// ProcessRequest는 설정에 따라 적절한 해시 처리를 수행
func (h *HashMaker) ProcessRequest() error {
	if h.config.UseZip {
		if h.config.ZipPath == "" {
			return fmt.Errorf("ZIP 모드에서는 ZIP 파일 경로가 필요합니다")
		}
		return h.ProcessZipFile(h.config.ZipPath)
	}

	if h.config.ZipFolder != "" && h.config.ZipName != "" {
		return h.CreateAndProcessZip()
	}

	return h.ProcessBasicHash()
}

// ProcessBasicHash는 기본 디렉토리 해싱을 처리
func (h *HashMaker) ProcessBasicHash() error {
	rootPath := h.config.StartPath
	if rootPath == "" {
		exePath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("실행 파일 경로 가져오기 실패: %w", err)
		}
		rootPath = filepath.Dir(exePath)
	}

	fmt.Printf("디렉토리 해시 계산 중: %s\n", rootPath)

	if err := h.calculateHashes(rootPath); err != nil {
		return fmt.Errorf("해시 계산 중 오류: %w", err)
	}

	sumFilePath := filepath.Join(rootPath, SumFileName)
	if err := h.writeSumFile(sumFilePath); err != nil {
		return fmt.Errorf("해시 파일 작성 중 오류: %w", err)
	}

	fmt.Printf("해시 파일 생성 완료: '%s'\n", sumFilePath)
	return nil
}

```
## services/hash_service.go
```go
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

```
## services/zip_service.go
```go
package services

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CreateAndProcessZip은 디렉토리를 ZIP으로 압축하고 해시를 생성
func (h *HashMaker) CreateAndProcessZip() error {
	outputPath := h.config.ZipOutputPath
	if outputPath == "" {
		outputPath = "."
	}

	zipFilePath := filepath.Join(outputPath, h.config.ZipName)
	if err := h.zipDirectory(h.config.ZipFolder, zipFilePath); err != nil {
		return fmt.Errorf("ZIP 파일 생성 실패: %w", err)
	}

	fmt.Printf("ZIP 파일 생성 완료: %s\n", zipFilePath)
	return h.ProcessZipFile(zipFilePath)
}

// zipDirectory는 디렉토리를 ZIP 파일로 압축
func (h *HashMaker) zipDirectory(sourceDir, zipFile string) error {
	// ZIP 파일 생성을 위한 디렉토리 생성
	if err := os.MkdirAll(filepath.Dir(zipFile), os.ModePerm); err != nil {
		return fmt.Errorf("ZIP 디렉토리 생성 실패: %w", err)
	}

	// ZIP 파일 생성
	zipfile, err := os.Create(zipFile)
	if err != nil {
		return fmt.Errorf("ZIP 파일 생성 실패: %w", err)
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	// 디렉토리 순회하며 파일 압축
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// ZIP 엔트리 헤더 생성
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("ZIP 헤더 생성 실패: %w", err)
		}

		// 상대 경로 계산
		relPath, err := filepath.Rel(sourceDir, path)
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

		// 파일 내용 복사
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("파일 열기 실패: %w", err)
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})
}

```
