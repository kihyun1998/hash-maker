package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kihyun1998/hash-maker/internal/domain/model"
	"github.com/kihyun1998/hash-maker/internal/domain/repository"
	"github.com/kihyun1998/hash-maker/internal/infrastructure/config"
	"github.com/kihyun1998/hash-maker/internal/infrastructure/fsys"
	"github.com/kihyun1998/hash-maker/internal/infrastructure/hashgen"
	"github.com/kihyun1998/hash-maker/internal/infrastructure/ziputil"
	"github.com/kihyun1998/hash-maker/internal/service/archiver"
	"github.com/kihyun1998/hash-maker/internal/service/hasher"
)

func main() {
	// 설정 초기화
	config, err := config.NewFlagConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "설정 초기화 실패: %v\n", err)
		os.Exit(1)
	}

	if err := config.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "설정 검증 실패: %v\n", err)
		os.Exit(1)
	}

	// 의존성 초기화
	hashGenerator := hashgen.NewSHA256Generator()
	fileSystem := fsys.NewLocalFileSystem()
	zipArchiver := ziputil.NewZipArchiver()

	// 서비스 초기화
	hashService := hasher.NewHashService(hashGenerator, fileSystem, config)
	archiveService := archiver.NewArchiveService(zipArchiver, hashGenerator, fileSystem, config)

	// 요청 처리
	if err := processRequest(config, hashService, archiveService); err != nil {
		fmt.Fprintf(os.Stderr, "오류 발생: %v\n", err)
		os.Exit(1)
	}
}

// processRequest는 설정에 따라 적절한 서비스를 호출
func processRequest(
	config repository.IConfigProvider,
	hashService *hasher.HashService,
	archiveService *archiver.ArchiveService,
) error {
	// ZIP 모드일 경우
	if config.IsHashMode() {
		result, err := archiveService.ProcessArchive(config.GetZipPath())
		if err != nil {
			return err
		}
		printHashResult(result)
		return nil
	}

	// ZIP 생성 모드일 경우
	if config.GetZipFolder() != "" && config.GetZipName() != "" {
		outputPath := config.GetOutputPath()
		zipPath := filepath.Join(outputPath, config.GetZipName())
		result, err := archiveService.CreateAndProcessArchive(config.GetZipFolder(), zipPath)
		if err != nil {
			return err
		}
		printHashResult(result)
		return nil
	}

	// 기본 해시 모드일 경우
	return hashService.ProcessDirectory(config.GetStartPath())
}

// printHashResult는 해시 처리 결과를 출력
func printHashResult(result *model.HashResult) {
	fmt.Printf("해시 추가 완료\n")
	fmt.Printf("- 원본 크기: %d 바이트\n", result.OriginalSize)
	fmt.Printf("- 해시 길이: %d 바이트\n", result.HashLength)
	fmt.Printf("- 최종 크기: %d 바이트\n", result.FinalSize)
}
