package main

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var startPath string
var rootPath string
var zipPath string
var useZip bool
var zipFolder string
var zipName string
var zipOutputPath string

func init() {
	flag.StringVar(&startPath, "startPath", "", "The valiable startPath is start root path")
	flag.StringVar(&zipPath, "zipPath", "", "The path to the zip file to be hashed")
	flag.BoolVar(&useZip, "zip", false, "Use zip file mode")
	flag.StringVar(&zipFolder, "zipfolder", "", "The folder to be zipped")
	flag.StringVar(&zipName, "zipname", "", "The name of the output zip file")
	flag.StringVar(&zipOutputPath, "zipoutput", "", "The output path for the zip file")
	flag.Parse()
}

const (
	sumFileName    = "hash_sum.txt"
	executableName = "hash-maker.exe"
)

func main() {
	if useZip {
		if zipPath == "" {
			fmt.Println("Error: -zipPath is required when -zip flag is used")
			return
		}

		err := processZipFile(zipPath)
		if err != nil {
			fmt.Printf("Error processing zip file: %v\n", err)
			return
		}
		fmt.Println("Zip file processed successfully.")
		return
	}

	if zipFolder != "" && zipName != "" {
		outputPath := zipOutputPath
		if outputPath == "" {
			outputPath = "."
		}
		zipFilePath := filepath.Join(outputPath, zipName)

		err := zipDirectory(zipFolder, zipFilePath)
		if err != nil {
			fmt.Printf("Error creating zip file: %v\n", err)
			return
		}
		fmt.Printf("Successfully created zip file: %s\n", zipFilePath)
		return
	}

	processBasicHash()
}

func processZipFile(zipFilePath string) error {
	file, err := os.OpenFile(zipFilePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	originalSize := fileInfo.Size()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Errorf("failed to calculate hash: %v", err)
	}

	hashSum := hash.Sum(nil)
	hashLength := len(hashSum)

	// 파일 포인터를 파일의 끝으로 이동합니다.
	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		return fmt.Errorf("failed to seek to end of file: %v", err)
	}

	// 해시 길이를 먼저 씁니다 (4바이트 정수).
	if err := binary.Write(file, binary.LittleEndian, int32(hashLength)); err != nil {
		return fmt.Errorf("failed to write hash length: %v", err)
	}

	// 해시를 파일의 끝에 씁니다.
	if _, err := file.Write(hashSum); err != nil {
		return fmt.Errorf("failed to write hash to file: %v", err)
	}

	fmt.Printf("Hash appended to file. Original size: %d bytes, Hash length: %d bytes, New size: %d bytes\n",
		originalSize, hashLength, originalSize+int64(hashLength)+4)

	return nil
}

func processBasicHash() {
	if startPath != "" {
		rootPath = filepath.Dir(startPath)
	} else {
		exePath, err := os.Executable()
		if err != nil {
			fmt.Printf("Error getting executable path: %v\n", err)
			return
		}

		rootPath = filepath.Dir(exePath)
	}

	fmt.Printf("Calculating hashes for directory: %s\n", rootPath)

	hashes, err := calculateHashes(rootPath)
	if err != nil {
		fmt.Printf("Error calculating hashes: %v\n", err)
		return
	}

	sumFilePath := filepath.Join(rootPath, sumFileName)
	err = writeSumFile(hashes, sumFilePath)
	if err != nil {
		fmt.Printf("Error writing sum file: %v\n", err)
		return
	}

	fmt.Printf("Hash sum file '%s' has been created successfully.\n", sumFilePath)
}
func zipDirectory(sourceDir, zipFile string) error {
	// ZIP 파일 생성 전 디렉토리 확인 및 생성
	zipDir := filepath.Dir(zipFile)
	if err := os.MkdirAll(zipDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory for zip file: %v", err)
	}

	// ZIP 파일 생성
	zipfile, err := os.Create(zipFile)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	// 소스 디렉토리를 순회하며 파일 압축
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// ZIP 파일 내 경로 생성
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(relPath)

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
