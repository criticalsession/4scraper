package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
	"github.com/schollz/progressbar/v3"
)

type DownloadableFile struct {
	Url      string
	FileName string
}

func main() {
	var threadUrl string
	fmt.Print("> Enter thread url: ")
	fmt.Scan(&threadUrl)

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

	c.OnHTML(".fileText > a", func(h *colly.HTMLElement) {
		files = append(files, DownloadableFile{
			FileName: h.Text,
			Url:      h.Request.AbsoluteURL(h.Attr("href")),
		})
	})

	c.Visit(threadUrl)

	if len(files) == 0 {
		fmt.Println("No files found to download.")
		return
	}

	bar := progressbar.Default(int64(len(files)))
	for _, file := range files {
		err := DownloadFile(file.Url, downloadDir, file.FileName)

		if err != nil {
			panic(err.Error())
		}

		bar.Add(1)
	}

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("-> Visiting", r.URL.String())
	// })

	fmt.Println("Download finished. Thank you for using 4scraper!")
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
