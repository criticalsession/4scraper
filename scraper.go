package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	fmt.Println("* 4scraper *")

	c := colly.NewCollector(
		colly.AllowedDomains("boards.4chan.org"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("-> Visiting", r.URL.String())
	})

	c.Visit("https://boards.4chan.org/x/thread/36739881")

	fmt.Println("Done.")
}
