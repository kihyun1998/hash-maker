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
