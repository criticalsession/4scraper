package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
)

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

	c.OnHTML(".fileText > a", func(h *colly.HTMLElement) {
		err := DownloadFile(h.Request.AbsoluteURL(h.Attr("href")), downloadDir, h.Text)
		if err != nil {
			panic(err.Error())
		}
	})

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("-> Visiting", r.URL.String())
	// })

	c.Visit(threadUrl)

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
