package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/common-nighthawk/go-figure"
	"github.com/criticalsession/4scraper/config"
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
	flags := config.ParseFlags()
	conf := config.ReadConfig()

	if flags.IsSearch {
		figure.NewFigure("4scraper", "rectangles", true).Print()
		fmt.Printf("%s [%s]\n", config.Version, "Search Mode")
		colorstring.Println("[dim]https://github.com/criticalsession/4scraper")
		fmt.Println("")

		threads, err := search.FindInBoard(flags.SearchBoard, flags.SearchTerms)
		if err != nil {
			logErr(err)
			return
		}

		if len(threads) > 0 {
			colorstring.Printf("Search results for [green]%s[default] in [green]/%s/[default]\n\n", strings.Join(flags.SearchTerms, ", "), flags.SearchBoard)
			for i, t := range threads {
				colorstring.Printf("[bold][green][%d][default] https://boards.4chan.org/%s/thread/%d (Images: %d)\n", i+1, flags.SearchBoard, t.No, t.Images+1)
				if t.Sub != "" {
					fmt.Printf("  %s\n", t.Sub)
				}

				if t.Com != "" {
					fmt.Printf("  %s\n", t.Com)
				}

				fmt.Println()
			}

			var downOption string
			colorstring.Print("> Enter [green][number][default] to download thread, [green]all[default] to download all, [green]q[default] to quit: ")
			fmt.Fscan(os.Stdin, &downOption)

			downOption = strings.ToLower(downOption)

			threadId, err := strconv.ParseInt(downOption, 10, 32)
			if err != nil {
				if downOption == "all" || downOption == "a" {
					for _, t := range threads {
						url := fmt.Sprintf("https://boards.4chan.org/%s/thread/%d", flags.SearchBoard, t.No)
						colorstring.Printf("\nDownloading thread [green]%d[default]...\n", t.No)
						downloadUrl(url, flags, conf)
					}
				} else {
					return
				}
			} else {
				threadId -= 1
				if threadId >= 0 && int(threadId) < len(threads) {
					colorstring.Printf("\nDownloading thread [green]%d[default]...\n", threads[threadId].No)
					url := fmt.Sprintf("https://boards.4chan.org/%s/thread/%d", flags.SearchBoard, threads[threadId].No)
					downloadUrl(url, flags, conf)
				} else {
					logErr(errors.New("invalid thread number"))
				}
			}
		} else {
			fmt.Printf("No results found for %s in /%s/\n\n", strings.Join(flags.SearchTerms, ", "), flags.SearchBoard)
			return
		}
	} else {
		url := flags.Url

		// head
		if !flags.Silent {
			figure.NewFigure("4scraper", "rectangles", true).Print()
			fmt.Printf("%s [%s]\n", config.Version, "Download Mode")
			fmt.Println("https://github.com/criticalsession/4scraper")
			fmt.Println("")

			if url == "" {
				fmt.Print("> Enter thread url: ")
				fmt.Fscan(os.Stdin, &url)
			} else {
				fmt.Println("Using arg url:", url)
			}
		}

		downloadUrl(url, flags, conf)
		if !flags.Silent {
			fmt.Println("\nDownload finished. Thank you for using 4scraper!")
		}
	}
}

func downloadUrl(url string, flags config.FlagSettings, conf config.Config) {
	// build download dir
	board, threadId, err := extractBoardAndThreadId(url)
	if err != nil {
		logErr(err)
		return
	}

	var downloadDir string
	if len(flags.OutDir) > 0 {
		downloadDir = flags.OutDir
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

		if !flags.Silent {
			bar.Add(1)
		}
	})

	c.Visit(url)

	if len(files) == 0 && !flags.Silent {
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
			if err := DownloadFile(file.Url, downloadDir, file.FileName, conf.UseOriginalFilename, flags.Silent, bar); err != nil {
				logErr(err)
			}
		} else {
			wg.Add(1)
			queue <- true
			go func(file DownloadableFile, i int) {
				defer wg.Done()
				if err := DownloadFile(file.Url, downloadDir, file.FileName, conf.UseOriginalFilename, flags.Silent, bar); err != nil {
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
