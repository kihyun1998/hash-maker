package models

// FileInfo는 파일의 해시 정보를 저장하는 구조체
type FileInfo struct {
	FileType string // 파일 타입
	DataHash string // 파일 내용의 해시값
	PathHash string // 파일 경로의 해시값
}
