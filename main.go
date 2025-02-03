// package main

// import (
// 	"flag"
// 	"fmt"
// 	"os"

// 	"github.com/kihyun1998/hash-maker/models"
// 	"github.com/kihyun1998/hash-maker/services"
// )

// func main() {
// 	// 명령줄 인자 파싱
// 	config := parseFlags()

// 	// HashMaker 인스턴스 생성
// 	hashMaker := services.NewHashMaker(config)

// 	// 해시 처리 실행
// 	if err := hashMaker.ProcessRequest(); err != nil {
// 		fmt.Fprintf(os.Stderr, "오류 발생: %v\n", err)
// 		os.Exit(1)
// 	}
// }

// // parseFlags는 명령줄 인자를 파싱하여 Config 구조체를 반환
// func parseFlags() models.Config {
// 	config := models.Config{}

// 	flag.StringVar(&config.StartPath, "startPath", "", "시작 경로 지정")
// 	flag.StringVar(&config.ZipPath, "zipPath", "", "ZIP 파일 경로 지정")
// 	flag.BoolVar(&config.UseZip, "zip", false, "ZIP 모드 사용")
// 	flag.StringVar(&config.ZipFolder, "zipfolder", "", "ZIP으로 만들 폴더 지정")
// 	flag.StringVar(&config.ZipName, "zipname", "", "생성할 ZIP 파일 이름")
// 	flag.StringVar(&config.ZipOutputPath, "zipoutput", "", "ZIP 파일 출력 경로")

// 	flag.Parse()

// 	return config
// }

// cmd/main.go
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
