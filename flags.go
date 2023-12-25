package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

const v = "v1.3" // unreleased

func ParseFlags() (bool, string, string) {
	var sil, help, ver bool
	var url, outDir string

	flag.BoolVar(&sil, "silent", false, "run without output, requires [url] arg")
	flag.BoolVar(&sil, "s", false, "run without output, requires [url] arg")

	flag.BoolVar(&help, "help", false, "show help")
	flag.BoolVar(&help, "h", false, "show help")

	flag.BoolVar(&ver, "version", false, "show version")
	flag.BoolVar(&ver, "v", false, "show version")

	flag.StringVar(&outDir, "output", "", "output directory")
	flag.StringVar(&outDir, "o", "", "output directory")

	flag.Usage = printHelp
	flag.Parse()

	url = flag.Arg(0)

	if url == "" {
		if help {
			printHelp()
			os.Exit(0)
		} else if ver {
			printVersion()
			os.Exit(0)
		}

		sil = false // silent cannot be used without a url
	}

	return sil, outDir, url
}

func printHelp() {
	n := "./4scraper.bin"
	if runtime.GOOS == "windows" {
		n = "4scraper.exe"
	}

	fmt.Printf(`Usage: %s [options] [URL]

Description:
  4scraper is a command line tool that quickly finds and downloads all images, 
  videos and gifs in a given 4chan thread.

Options:
  -h, --help         Show this help message and exit
  -o, --output       Specify output directory for downloaded files
  -s, --silent       Run in silent mode (no output), requires URL
  -v, --version      Show version number and exit

Arguments:
  URL                Full 4chan thread URL to scrape

If no URL is provided, --silent flag will be ignored and 4scraper will ask you 
to enter a thread URL.

A valid URL looks like this:
  https://boards.4chan.org/[BOARD]/thread/[THREAD_ID]
  Example: https://boards.4chan.org/g/thread/76759434 (sfw)

Examples:
  # Default execution with prompt
  %s

  # Scrape thread without user input
  %s https://boards.4chan.org/g/thread/76759434
  
  # Scrape thread without user input or output
  %s --silent https://boards.4chan.org/g/thread/76759434
  
  # Scrape thread and store in custom directory
  %s --output=downloads/battlestations https://boards.4chan.org/g/thread/76759434

  Source:
  https://github.com/criticalsession/4scraper
`, n, n, n, n, n)
}

func printVersion() {
	fmt.Println("4scraper", v)
}
