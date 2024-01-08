# 4scraper

## 4chan scraper CLI tool written in Go that downloads all images, videos and gifs from a 4chan thread (no setup required)!

[![GitHub release (with filter)](https://img.shields.io/github/v/release/criticalsession/4scraper)](https://github.com/criticalsession/4scraper/releases)
[![GitHub issues](https://img.shields.io/github/issues/criticalsession/4scraper)](https://github.com/criticalsession/4scraper/issues)
[![Go Report Card](https://goreportcard.com/badge/github.com/criticalsession/4scraper)](https://goreportcard.com/report/github.com/criticalsession/4scraper)
[![GitHub License](https://img.shields.io/github/license/criticalsession/4scraper)](https://github.com/criticalsession/4scraper/blob/main/LICENSE)
[![X (formerly Twitter) Follow](https://img.shields.io/twitter/follow/criticalsession)](https://twitter.com/criticalsession)

<p align="center">
  <img src="https://github.com/criticalsession/4scraper/blob/main/docs/4scraper_search.gif?raw=true" width="608" />
</p>

<p align="center"><b>Version: 1.4</b> - Thread Search Edition</p>

**4scraper** is an open source command line tool written in Go that quickly finds and downloads all images, videos and gifs in a given thread. No setup or installation required, and no fluff.

1. Run **4scraper**
2. Enter/paste the thread you'd like to download from and press `Enter`
3. Wait for the download to finish
4. Profit

## Contents

- [1. How to install](#1-how-to-install)
- [2. How to use](#2-how-to-use)
    - [2.1 Args and Flags](#21-args-and-flags)
    - [2.2 Configuration](#22-configuration)
- [3. How it works](#3-how-it-works)
- [4. Found a bug? Have suggestions?](#4-found-a-bug-have-suggestions)
- [5. Known issues](#5-known-issues)
- [6. Like 4scraper?](#6-like-4scraper)

## 1. How to install

Head on over to the <a href="https://github.com/criticalsession/4scraper/releases">releases tab</a> and pick the version you need. There's no setup or installation required; simply run the downloaded `exe` or `bin` file and you're good to go. Or you can download the code and build it yourself.

**Optional:** Create a config file as outlined in [2.2 Configuration](#22-configuration-v12)

## 2. How to use

Click on the thumbnail below to watch the **shortest video tutorial** of your life (it's a YouTube link).

<a href="https://www.youtube.com/watch?v=2cLXpOMFGdc" target="_blank">
    <img src="https://img.youtube.com/vi/2cLXpOMFGdc/0.jpg" alt="Watch the tutorial" width="240"/>
</a>

### 2.1 Args and Flags

Added in v1.1, **4scraper** can now be executed in silent mode by setting the flag and passing a URL as arg. Here's a full description of flags available. More information is available at `4scraper.exe --help`. For brevity I'm using `4scraper.exe` but on linux it would be `./4scraper.bin`:

```
Usage: 
  4scraper.exe [options] [URL]

  [URL]                
    Full 4chan thread URL to download files from

OPTIONS:
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
```
#### Additional Notes

- If no `URL` is provided, the `--silent` flag will be ignored and **4scraper** will ask you to enter a thread URL as if you executed without flags.
- If an `output` directory is specified, it will override the `downloads/board/threadid` directory structure and the `BoardDir` `ThreadDir` options set in config will be ignored.
- When using `--find`, the first argument should be the board code, followed by search terms. Search terms must all be found for the thread to be returned.

### 2.2 Configuration

As of v1.2 you can add a configuration file to setup basic settings. This is entirely optional and the software will run even if there's no config file set. If you'd like to create and setup a config file read on.

1. Create a file named `4scraper.config` in the same directory as the **4scraper** executable
2. Copy and paste the following inside it and save
```
# 4scraper config file
BoardDir = true
ThreadDir = true
UseOriginalFilename = true
ParallelDownload = true
Extensions = "jpeg,jpg,png,gif,webp,bpm,tiff,svg,mp4,mov,avi,webm,flv"
```
3. Lines starting with `#` are ignored
4. All settings should be in the format `[key] = [value]`
5. These are the settings you can adjust:
      1. `BoardDir`: if `true` a directory with the board code will be created in the `downloads` directory for organization (e.g. `downloads/g/`)
      2. `ThreadDir`: if `true` a directory with the thread id will be created in the `downloads` or `board` directory for organization (e.g. `downloads/g/4568995/`)
      3. `UseOriginalFilename`: if `false` a new unique filename will be generated (using `UUID`)
      4. `ParallelDownload`: if `true` downloads files concurrently up to a maximum of `20` concurrent threads
      5. `Extensions`: any file type that isn't in the list won't be downloaded
6. If both `BoardDir` and `ThreadDir` are turned off, all downloaded files will go in the `downloads/` directory

If no config file is created, the setting defaults are as shown here.

## 3. How it works

**4scraper** was developed as an exercise in getting my hands dirty with Go, so there's nothing wild going on behind the scenes. Still, this is how it works, in case anyone's interested.

1. Once a thread URL is provided, the `board` and `threadId` are extracted for later, and an indeterminate [ProgressBar](https://github.com/schollz/progressbar/) is initiated and shown on screen
2. Next, [Colly](https://github.com/gocolly/colly) is used to scrape and retrieve all file URLs and filenames in the thread; files are found by looking for `.filetext > a`
3. All files found are stored in a `[]DownloadableFile slice` for later and the `ProgressBar` is updated to reflect the total files to download
4. We create directories in this structure: `downloads/<board>/<threadId>` to hold all downloaded files
5. For each `DownloadableFile` found we first check if the filename already exists (and append a random number to the filename if it does) then we download it from 4chan
    - if `ParallelDownloads` is `true`, the download code is called in a go coroutine with a queue that manages maximum threads

## 4. Found a bug? Have suggestions?

Feel free to use the Issues tab above (or [click here](https://github.com/criticalsession/4scraper/issues)) if you've found bugs, have problems running **4scraper**, have suggestions for improvements or general tips on how I can make the Go code better.

## 5. Known issues

- None

## 6. Like 4scraper?

If you're feeling generous, buy me a beer! - https://www.buymeacoffee.com/criticalsession 🍺❤️
