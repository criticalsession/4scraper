package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/common-nighthawk/go-figure"
	"github.com/criticalsession/4scraper/search"
	"github.com/gocolly/colly"
	"github.com/k0kubun/go-ansi"
	"github.com/mitchellh/colorstring"
	"github.com/schollz/progressbar/v3"
)

type DownloadableFile struct {
	Url      string
	FileName string
}

func main() {
	// setup
	flags := ParseFlags()
	silent := flags.Silent
	outDir := flags.OutDir
	url := flags.Url

	if flags.IsSearch {
		threads, err := search.FindInBoard(flags.SearchBoard, flags.SearchTerms)
		if err != nil {
			logErr(err)
			return
		}

		for _, t := range threads {
			fmt.Printf("https://boards.4chan.org/%s/thread/%d\n", flags.SearchBoard, t.No)
		}
	} else {
		// config
		conf := ReadConfig()

		// head
		if !silent {
			figure.NewFigure("4scraper", "rectangles", true).Print()
			fmt.Println(v)
			fmt.Println("https://github.com/criticalsession/4scraper")
			fmt.Println("")

			if url == "" {
				fmt.Print("> Enter thread url: ")
				fmt.Fscan(os.Stdin, &url)
			} else {
				fmt.Println("Using arg url:", url)
			}
		}

		// build download dir
		board, threadId, err := extractBoardAndThreadId(url)
		if err != nil {
			logErr(err)
			return
		}

		var downloadDir string
		if len(outDir) > 0 {
			downloadDir = outDir
		} else {
			downloadDir = "downloads"
			if conf.BoardDir {
				downloadDir += "/" + board
			}

			if conf.ThreadDir {
				downloadDir += "/" + threadId
			}
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
			addLinkToDownloadableFiles(h, &conf.Extensions, &files)

			if !silent {
				bar.Add(1)
			}
		})

		c.Visit(url)

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

		var wg sync.WaitGroup
		queue := make(chan bool, min(20, len(files)))
		errs := []error{}

		for i, file := range files {
			if !conf.ParallelDownload {
				if err := DownloadFile(file.Url, downloadDir, file.FileName, conf.UseOriginalFilename, silent, bar); err != nil {
					logErr(err)
				}
			} else {
				wg.Add(1)
				queue <- true
				go func(file DownloadableFile, i int) {
					defer wg.Done()
					if err := DownloadFile(file.Url, downloadDir, file.FileName, conf.UseOriginalFilename, silent, bar); err != nil {
						errs = append(errs, err)
					}

					<-queue
				}(file, i)
			}
		}

		if conf.ParallelDownload {
			wg.Wait()

			for _, err := range errs {
				logErr(err)
			}
		}

		if !silent {
			fmt.Println("\nDownload finished. Thank you for using 4scraper!")
		}
	}
}

func extractBoardAndThreadId(urlStr string) (string, string, error) {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return "", "", err
	}

	if !strings.Contains(urlStr, "boards.4chan.org") {
		return "", "", errors.New("not a 4chan URL")
	}

	urlStr = strings.TrimPrefix(parsed.Path, "boards.4chan.org")
	urlStr = strings.TrimPrefix(urlStr, "/")

	if len(strings.Split(urlStr, "/")) < 3 {
		return "", "", errors.New("not a valid thread URL")
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
		logErr(err)
		return
	}

	urlSpl := strings.Split(parsed.Path, "/")
	ext := urlSpl[len(urlSpl)-1]
	ext = strings.Split(ext, ".")[len(strings.Split(ext, "."))-1]

	if slices.Contains(*extensions, ext) {
		// add file to files to download
		*files = append(*files, DownloadableFile{
			FileName: h.Text,
			Url:      fileUrl,
		})
	}
}

func logErr(err error) {
	colorstring.Printf("\n[red]ERROR: %s\n", err.Error())
}
