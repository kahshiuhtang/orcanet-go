package store

import (
	"log"
	"os"
	"time"
)

type FileInfo struct {
	IsDir   bool
	ModTime time.Time
	Name    string
	Size    int64
}

func GetAllFiles() []FileInfo {
	files, err := os.ReadDir("files")
	if err != nil {
		log.Fatal(err)
	}
	fileNames := make([]FileInfo, 1)
	for _, file := range files {
		fileInfo, err := os.Stat("files/" + file.Name())
		if err == nil {
			fileNames = append(fileNames, FileInfo{IsDir: fileInfo.IsDir(), ModTime: fileInfo.ModTime(), Name: fileInfo.Name(), Size: fileInfo.Size()})
		}
	}
	return fileNames
}
