package repository

// IPlatformProvider는 플랫폼 관련 정보를 제공하는 인터페이스
type IPlatformProvider interface {
	GetExecutableName() string
	GetCurrentExecutablePath() (string, error)
}
