package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"slices"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/gocolly/colly"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

type DownloadableFile struct {
	Url      string
	FileName string
}

func main() {
	// setup
	silent, threadUrl := ParseFlags()
	config := ReadConfig()

	// head
	if !silent {
		figure.NewFigure("4scraper", "rectangles", true).Print()
		fmt.Println("v1.2")
		fmt.Println("https://github.com/criticalsession/4scraper")
		fmt.Println("")

		if threadUrl == "" {
			fmt.Print("> Enter thread url: ")
			fmt.Fscan(os.Stdin, &threadUrl)
		} else {
			fmt.Println("Using arg:", threadUrl)
		}
	}

	// build download dir
	board, threadId, err := extractBoardAndThreadId(threadUrl)
	if err != nil {
		log.Fatalln("ERROR:", err.Error())
		return
	}

	downloadDir := "downloads"
	if config.BoardDir {
		downloadDir += "/" + board
	}

	if config.ThreadDir {
		downloadDir += "/" + threadId
	}

	c := colly.NewCollector(
		colly.AllowedDomains("boards.4chan.org", "i.4cdn.org"),
	)

	files := []DownloadableFile{}

	// find files
	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("[cyan][1/2][reset] Finding files:"),
	)

	c.OnHTML(".fileText > a", func(h *colly.HTMLElement) {
		addLinkToDownloadableFiles(h, &config.Extensions, &files)

		if !silent {
			bar.Add(1)
		}
	})

	c.Visit(threadUrl)

	if len(files) == 0 && !silent {
		fmt.Println("No files found to download.")
		return
	}

	// download files
	bar.Reset()
	bar = progressbar.NewOptions(len(files),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionShowCount(),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetDescription("[cyan][2/2][reset] Downloading:"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	for _, file := range files {
		err := DownloadFile(file.Url, downloadDir, file.FileName, config.UseOriginalFilename)

		if err != nil {
			log.Fatalln("ERROR:", err.Error())
		}

		if !silent {
			bar.Add(1)
		}
	}

	if !silent {
		fmt.Println("\nDownload finished. Thank you for using 4scraper!")
	}
}

func extractBoardAndThreadId(urlStr string) (string, string, error) {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return "", "", err
	}

	if !strings.Contains(urlStr, "boards.4chan.org") {
		return "", "", errors.New("Not a 4chan URL")
	}

	urlStr = strings.TrimPrefix(parsed.Path, "boards.4chan.org")
	urlStr = strings.TrimPrefix(urlStr, "/")

	if len(strings.Split(urlStr, "/")) < 3 {
		return "", "", errors.New("Not a valid thread URL")
	}

	board := strings.Split(urlStr, "/")[0]
	threadId := strings.Split(urlStr, "/")[2]

	return board, threadId, nil
}

func addLinkToDownloadableFiles(h *colly.HTMLElement, extensions *[]string, files *[]DownloadableFile) {
	fileUrl := h.Request.AbsoluteURL(h.Attr("href"))

	// check if extension is downloadable
	parsed, err := url.Parse(fileUrl)
	if err != nil {
		log.Fatalln("ERROR:", err.Error())
	}

	fileUrlSplit := strings.Split(parsed.Path, "/")
	ext := fileUrlSplit[len(fileUrlSplit)-1]
	ext = strings.Split(ext, ".")[len(strings.Split(ext, "."))-1]

	if slices.Contains(*extensions, ext) {
		// add file to files to download
		*files = append(*files, DownloadableFile{
			FileName: h.Text,
			Url:      fileUrl,
		})
	}
}
