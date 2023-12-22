package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	fmt.Println("* 4scraper *")

	c := colly.NewCollector(
		colly.AllowedDomains("boards.4chan.org", "i.4cdn.org"),
	)

	c.OnHTML(".fileText > a", func(h *colly.HTMLElement) {
		fmt.Printf("--> Downloading: %v\n", h.Text)
		err := DownloadFile(h.Request.AbsoluteURL(h.Attr("href")), h.Text)
		if err != nil {
			panic(err.Error())
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("-> Visiting", r.URL.String())
	})

	c.Visit("https://boards.4chan.org/x/thread/36750486")

	fmt.Println("Done.")
}
