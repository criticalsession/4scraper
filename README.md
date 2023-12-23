# 4scraper

## 4chan scraper CLI tool written in Go that downloads all images, videos and gifs from a 4chan thread (no setup required)!

<p align="center">
  <img src="https://github.com/criticalsession/4scraper/blob/feat/config/docs/scraper-main.1_2.png?raw=true" />
</p>

<p align="center"><b>Version: 1.2</b></p>

**4scraper** is an open source command line tool written in Go that quickly finds and downloads all images, videos and gifs in a given thread. No setup or installation required, and no fluff.

* Run **4scraper**
* Enter/paste the thread you'd like to download from and press `Enter`
* Wait for the download to finish
* Profit

## Contents

- [1. How to install](#1-how-to-install)
- [2. How to use](#2-how-to-use)
    - [2.1 Args and Flags](#21-args-and-flags-v11)
    - [2.2 Configuration](#22-configuration-v12)
- [3. How it works](#3-how-it-works)
- [4. Found a bug? Have suggestions?](#4-found-a-bug-have-suggestions)
- [5. Things I want to add](#5-things-i-want-to-add)
- [6. Like 4scraper?](#6-like-4scraper)

## 1. How to install

Head on over to the <a href="https://github.com/criticalsession/4scraper/releases">releases tab</a> and pick the version you need. There's no setup or installation required; simply run the downloaded `exe` or `bin` file and you're good to go. Or you can download the code and build it yourself.

**Optional:** Create a config file as outlined in [2.2 Configuration](#22-configuration-v12)

## 2. How to use

Click on the thumbnail below to watch the **shortest video tutorial** of your life (it's a YouTube link).

<a href="https://www.youtube.com/watch?v=2cLXpOMFGdc" target="_blank">
    <img src="https://img.youtube.com/vi/2cLXpOMFGdc/0.jpg" alt="Watch the tutorial" width="240"/>
</a>

### 2.1 Args and Flags [v1.1]

Added in v1.1, **4scraper** can now be executed in silent mode by setting the flag and passing a URL as arg. Here's a full description of flags available. More information is available at `4scraper.exe --help`. For brevity I'm using `4scraper.exe` but on linux it would be `./4scraper.bin`:

```
Usage: 4scraper.exe [options] [URL]

Options:
  -h, --help         Show this help message and exit
  -s, --silent       Run in silent mode (no output), requires URL
  -v, --version      Show version number and exit

Arguments:
  URL                Full 4chan thread URL to scrape
```

In short:
- `4scraper.exe` will execute as normal; it will prompt you for a `thread URL` and display verbose progress information
- `4scraper.exe https://boards.4chan.org/g/thread/76759434` will execute normally (verbose) but won't ask you for a `thread URL`
- `4scraper.exe --silent https://boards.4chan.org/g/thread/76759434` will execute using the `thread URL` provided and won't display verbose output, progress, etc.

**NOTE:** If no `URL` is provided, the `--silent` flag will be ignored and **4scraper** will ask you to enter a thread URL as if you executed without flags.

### 2.2 Configuration [v1.2]

As of v1.2 you can add a configuration file to setup basic settings. This is entirely optional and the software will run even if there's no config file set. If you'd like to create and setup a config file read on.

1. Create a file named `4scraper.config` in the same directory as the **4scraper** executable
2. Copy and paste the following inside it and save
```
# 4scraper config file
BoardDir = true
ThreadDir = true
UseOriginalFilename = true
Extensions = "jpeg,jpg,png,gif,webp,bpm,tiff,svg,mp4,mov,avi,webm,flv"
```
3. Lines starting with `#` are ignored
4. All settings should be in the format `[key] = [value]`
5. These are the settings you can adjust:
      1. `BoardDir`: if `true` a directory with the board code will be created in the `downloads` directory for organization (e.g. `downloads/g/`)
      2. `ThreadDir`: if `true` a directory with the thread id will be created in the `downloads` or `board` directory for organization (e.g. `downloads/g/4568995/`)
      3. `Extensions`: any file that isn't in the list won't be downloaded
      4. `UseOriginalFilename`: if `false` a new unique filename will be generated (using `UUID`)
6. If both `BoardDir` and `ThreadDir` are turned off, all downloaded files will go in the `downloads/` directory

## 3. How it works

**4scraper** was developed as an exercise in getting my hands dirty with Go, so there's nothing wild going on behind the scenes. Still, this is how it works, in case anyone's interested.

1. Once a thread URL is provided, the `board` and `threadId` are extracted for later, and an indeterminate [ProgressBar](https://github.com/schollz/progressbar/) is initiated and shown on screen
2. Next, [Colly](https://github.com/gocolly/colly) is used to scrape and retrieve all file URLs and filenames in the thread; files are found by looking for `.filetext > a`
3. All files found are stored in a `[]DownloadableFile slice` for later and the `ProgressBar` is updated to reflect the total files to download
4. We create directories in this structure: `downloads/<board>/<threadId>` to hold all downloaded files
5. For each `DownloadableFile` found we first check if the filename already exists (and append a random number to the filename if it does) then we download it from 4chan

## 4. Found a bug? Have suggestions?

Feel free to use the Issues tab above (or [click here](https://github.com/criticalsession/4scraper/issues)) if you've found bugs, have problems running **4scraper**, have suggestions for improvements or general tips on how I can make the Go code better.

## 5. Things I want to add

- [ ] Optional config file for basic settings (folder organization, types of files to download)
- [x] Args to skip input with zero feedback to allow for automation

## 6. Like 4scraper?

If you're feeling generous, buy me a beer! - https://www.buymeacoffee.com/criticalsession 🍺❤️