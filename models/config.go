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
