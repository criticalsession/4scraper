package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/mrz1836/go-sanitize"
)

func DownloadFile(url, dir, filename string, useOriginalFilename bool) error {
	if !DirExists(dir) {
		os.MkdirAll(dir, os.ModePerm)
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if useOriginalFilename {
		filename = strings.Replace(filename, "(...)", "", -1) // remove (...) from the filename
	} else {
		filename = genUniqueFilename(filename)
	}

	baseName, ext := SplitFilename(filename)
	baseName = sanitize.AlphaNumeric(baseName, false)
	filename = fmt.Sprintf("%s%s", baseName, ext)

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

	return err
}

func tryGenNewFilename(dir, filename string) string {
	uniqueId := rand.Intn(100000)

	baseName, ext := SplitFilename(filename)

	newFilename := fmt.Sprintf("%s.%s%s", baseName, fmt.Sprint(uniqueId), ext)
	if FileExists(dir, newFilename) {
		return tryGenNewFilename(dir, filename)
	}

	return newFilename
}

func genUniqueFilename(filename string) string {
	uniqueId := strings.Replace(uuid.NewString(), "-", "", -1)
	ext := GetExtension(filename)

	return fmt.Sprintf("%s%s", uniqueId, ext)
}
