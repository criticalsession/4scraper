package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	fmt.Println("* 4scraper *")

	var threadUrl string
	fmt.Print("> Thread url: ")
	fmt.Scan(&threadUrl)

	board, threadId, err := extractBoardAndThreadId(threadUrl)
	if err != nil {
		panic(err)
	}

	downloadDir := fmt.Sprintf("downloads/%s/%s", board, threadId)

	c := colly.NewCollector(
		colly.AllowedDomains("boards.4chan.org", "i.4cdn.org"),
	)

	c.OnHTML(".fileText > a", func(h *colly.HTMLElement) {
		fmt.Printf("--> Found file: %v\n", h.Text)
		err := DownloadFile(h.Request.AbsoluteURL(h.Attr("href")), downloadDir, h.Text)
		if err != nil {
			panic(err.Error())
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("-> Visiting", r.URL.String())
	})

	c.Visit(threadUrl)

	fmt.Println("Done.")
}

func extractBoardAndThreadId(urlStr string) (string, string, error) {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return "", "", err
	}

	urlStr = strings.TrimPrefix(parsed.Path, "boards.4chan.org")
	urlStr = strings.TrimPrefix(urlStr, "/")

	board := strings.Split(urlStr, "/")[0]
	threadId := strings.Split(urlStr, "/")[2]

	return board, threadId, nil
}
