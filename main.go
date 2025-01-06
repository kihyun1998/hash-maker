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
