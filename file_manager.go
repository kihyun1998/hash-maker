package main

import (
	"fmt"
	"os"
	"sort"
)

func writeSumFile(hashes map[string]fileInfo, sumFilePath string) error {
	file, err := os.Create(sumFilePath)
	if err != nil {
		return fmt.Errorf("failed to create sum file: %v", err)
	}
	defer file.Close()

	var paths []string
	for path := range hashes {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for _, path := range paths {
		info := hashes[path]
		_, err := fmt.Fprintf(file, "%s;%s;%s\n", info.fileType, info.pathHash, info.dataHash)
		if err != nil {
			return fmt.Errorf("failed to write to sum file: %v", err)
		}
	}

	return nil
}
