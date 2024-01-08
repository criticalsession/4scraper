package config

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

const Version = "v1.4"

type FlagSettings struct {
	Silent      bool
	OutDir      string
	Url         string
	IsSearch    bool
	SearchTerms []string
	SearchBoard string
}

func ParseFlags() FlagSettings {
	var sil, help, ver, isSearch bool
	var url, outDir, searchBoard string
	var searchTerms []string

	flag.BoolVar(&sil, "silent", false, "run without output, requires [url] arg")
	flag.BoolVar(&sil, "s", false, "run without output, requires [url] arg")

	flag.BoolVar(&help, "help", false, "show help")
	flag.BoolVar(&help, "h", false, "show help")

	flag.BoolVar(&ver, "version", false, "show version")
	flag.BoolVar(&ver, "v", false, "show version")

	flag.StringVar(&outDir, "output", "", "output directory")
	flag.StringVar(&outDir, "o", "", "output directory")

	flag.BoolVar(&isSearch, "find", false, "search threads")
	flag.BoolVar(&isSearch, "f", false, "search threads")

	flag.Usage = printHelp
	flag.Parse()

	if !isSearch {
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

		return FlagSettings{
			Silent:   sil,
			OutDir:   outDir,
			Url:      url,
			IsSearch: false,
		}
	} else {
		if len(flag.Args()) < 2 {
			printHelp()
			os.Exit(0)
		}

		searchBoard = flag.Arg(0)
		searchTerms = flag.Args()[1:]

		return FlagSettings{
			IsSearch:    true,
			SearchTerms: searchTerms,
			SearchBoard: searchBoard,
		}
	}
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
  -h, --help         
    Show this help message and exit
	  
  -v, --version      
    Show version number and exit
  
  -o, --output [DIRECTORY]
    Specify output directory for downloaded files
  
  -s, --silent       
    Run in silent mode (no output), requires URL
  
  -f, --find [BOARD] [KEYWORDS]
    Search for threads in the specified board that match the given keywords
    [BOARD] is the name of the 4chan board (e.g., 'g')
    [KEYWORDS] are the terms to search for (e.g., 'linux desktop')  

Arguments:
  [URL]
    Full 4chan thread URL to download files from

If no URL is provided, --silent flag will be ignored and 4scraper will ask you 
to enter a thread URL.

A valid URL looks like this:
  https://boards.4chan.org/[BOARD]/thread/[THREAD_ID]
  Example: https://boards.4chan.org/g/thread/76759434 (sfw)

Examples:
  # Default execution with prompt
  %s

  # Scrape thread without user input
  %s [URL]
  
  # Scrape thread without user input or output
  %s --silent [URL]
  
  # Scrape thread and store in custom directory
  # e.g. will store all downloads in 'downloads/battlestations'
  %s --output=downloads/battlestations [URL]

  # Find threads
  # e.g. will search for threads in /g/ that contain 'linux' AND 'desktop'
  %s --find g linux desktop

  Source:
    https://github.com/criticalsession/4scraper
`, n, n, n, n, n, n)
}

func printVersion() {
	fmt.Println("4scraper", Version)
}
