package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
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
	silent, threadUrl := ParseFlags()

	if !silent {
		figure.NewFigure("4scraper", "rectangles", true).Print()
		fmt.Println("https://github.com/criticalsession/4scraper")
		fmt.Println("")

		if threadUrl == "" {
			fmt.Print("> Enter thread url: ")
			fmt.Fscan(os.Stdin, &threadUrl)
		} else {
			fmt.Println("Using arg:", threadUrl)
		}
	}

	board, threadId, err := extractBoardAndThreadId(threadUrl)
	if err != nil {
		log.Fatalln("ERROR:", err.Error())
		return
	}

	downloadDir := fmt.Sprintf("downloads/%s/%s", board, threadId)

	c := colly.NewCollector(
		colly.AllowedDomains("boards.4chan.org", "i.4cdn.org"),
	)

	files := []DownloadableFile{}

	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("[cyan][1/2][reset] Finding files:"),
	)

	c.OnHTML(".fileText > a", func(h *colly.HTMLElement) {
		files = append(files, DownloadableFile{
			FileName: h.Text,
			Url:      h.Request.AbsoluteURL(h.Attr("href")),
		})

		if !silent {
			bar.Add(1)
		}
	})

	c.Visit(threadUrl)

	if len(files) == 0 && !silent {
		fmt.Println("No files found to download.")
		return
	}

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
		err := DownloadFile(file.Url, downloadDir, file.FileName)

		if err != nil {
			panic(err.Error())
		}

		if !silent {
			bar.Add(1)
		}
	}

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("-> Visiting", r.URL.String())
	// })

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
