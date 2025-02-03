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
├── cmd/
│   └── main.go                 # 애플리케이션 진입점
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
