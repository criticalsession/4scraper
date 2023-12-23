package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

func ParseFlags() (bool, string) {
	var silent, help, version bool
	var url string

	flag.BoolVar(&silent, "silent", false, "run without output, requires [url] arg")
	flag.BoolVar(&silent, "s", false, "run without output, requires [url] arg")

	flag.BoolVar(&help, "help", false, "show help")
	flag.BoolVar(&help, "h", false, "show help")

	flag.BoolVar(&version, "version", false, "show version")
	flag.BoolVar(&version, "v", false, "show version")

	flag.Usage = printHelp
	flag.Parse()

	url = flag.Arg(0)

	if url == "" {
		if help {
			printHelp()
			os.Exit(0)
		} else if version {
			printVersion()
			os.Exit(0)
		}

		silent = false // silent cannot be used without a url
	}

	return silent, url
}

func printHelp() {
	toolName := "4scraper"
	if runtime.GOOS == "windows" {
		toolName += ".exe"
	}

	fmt.Printf(`Usage: %s [options] [URL]

Description:
  4scraper is a command line tool that quickly finds and downloads all images, 
  videos and gifs in a given 4chan thread.

Options:
  -h, --help         Show this help message and exit
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

  # Scrape thread without user input")
  %s https://boards.4chan.org/g/thread/76759434
  
  # Scrape thread without user input or output")
  %s --silent https://boards.4chan.org/g/thread/76759434

Source:
  https://github.com/criticalsession/4scraper
`, toolName, toolName, toolName, toolName)
}

func printVersion() {
	fmt.Println("4scraper version 1.1")
}
