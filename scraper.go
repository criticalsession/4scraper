package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	board := "x"
	threadId := "36750486"

	threadUrl := fmt.Sprintf("https://boards.4chan.org/%s/thread/%s", board, threadId)
	downloadDir := fmt.Sprintf("downloads/%s/%s", board, threadId)

	fmt.Println("* 4scraper *")

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
