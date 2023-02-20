package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ToJpeg converts an image to jpeg
func ToJpeg(imageBytes []byte, fileName string) ([]byte, error) {
	contentType := http.DetectContentType(imageBytes)
	switch contentType {
	case "image/png":
		img, _, err := image.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return nil, err
		}
		out, _ := os.Create(fileName + ".jpeg")
		defer out.Close()
		var opts jpeg.Options
		opts.Quality = 100
		err = jpeg.Encode(out, img, &opts)
		if err != nil {
			log.Println(err)
		}
		buf := new(bytes.Buffer)
		return buf.Bytes(), nil
	case "image/jpeg":
		log.Println("work in progress")
		//	ffmeg.
	case "image/avif":
		log.Println("welcome")
	}
	return nil, fmt.Errorf("unable to convert %#v to jpeg", contentType)
}

// to check files present in given file path
func FilePathToCheck(filepaths string) ([]string, error) {
	var files []string
	err := filepath.Walk(filepaths, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func main() {
	var filePaths string
	fmt.Println("Enter the folder path where the file will be placed: ")
	fmt.Scanln(&filePaths)
	files, err := FilePathToCheck(filePaths)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(files); i++ {
		file, err := os.Open(files[i])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()
		fileInfo, _ := file.Stat()
		var size int64 = fileInfo.Size()
		bytes := make([]byte, size)
		buffer := bufio.NewReader(file)
		io.ReadFull(buffer, bytes)
		fileName := strings.Trim(fileInfo.Name(), ".png")
		ToJpeg(bytes, fileName)
	}
}
