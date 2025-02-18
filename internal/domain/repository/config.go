package repository

// IConfigProvider는 설정 관리를 위한 인터페이스
type IConfigProvider interface {
	GetStartPath() string
	GetZipPath() string
	IsHashMode() bool
	GetZipFolder() string
	GetZipName() string
	GetOutputPath() string
}
