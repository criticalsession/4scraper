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
		fmt.Printf("--> Text: %v, URL: %v\n", h.Text, h.Request.AbsoluteURL(h.Attr("href")))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("-> Visiting", r.URL.String())
	})

	c.Visit("https://boards.4chan.org/x/thread/36739881")

	fmt.Println("Done.")
}
