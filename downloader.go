package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadFile(url, dir, filename string) error {
	if !DirExists(dir) {
		os.MkdirAll(dir, os.ModePerm)
		fmt.Println("---> Created directory:", dir)
	}

	fmt.Printf("---> Downloading... ")

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if FileExists(dir, filename) {
		filename = tryGenNewFilename(dir, filename)
	}

	fullFilename := filepath.Join(dir, filename)

	file, err := os.Create(fullFilename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	fmt.Println(" done.")

	return err
}

func tryGenNewFilename(dir, filename string) string {
	uniqueId := rand.Intn(100000)

	ext := filepath.Ext(filename)
	baseName := filename[:len(filename)-len(ext)]

	newFilename := fmt.Sprintf("%s.%s%s", baseName, fmt.Sprint(uniqueId), ext)
	if FileExists(dir, newFilename) {
		return tryGenNewFilename(dir, filename)
	}

	return newFilename
}
