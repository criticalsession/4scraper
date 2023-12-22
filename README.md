# 4scraper

## 4chan scraper CLI tool written in Go that downloads all images, videos and gifs from a 4chan thread (no setup required)!

<p align="center">
  <img src="https://github.com/criticalsession/4scraper/blob/main/docs/scraper-main.png?raw=true" />
</p>

**4scraper** is an open source command line tool written in Go that quickly finds and downloads all images, videos and gifs in a given thread. No setup or installation required, and no fluff.

* Run **4scraper**
* Enter/paste the thread you'd like to download from and press `Enter`
* Wait for the download to finish
* Profit

## How to install

Head on over to the <a href="https://github.com/criticalsession/4scraper/releases">releases tab</a> and pick the version you need. There's no setup or installation required; simply run the downloaded `exe` or `sh` file and you're good to go.

## How to use

Click on the thumbnail below to watch the **shortest video tutorial** of your life (it's a YouTube link).

<a href="https://www.youtube.com/watch?v=2cLXpOMFGdc" target="_blank">
    <img src="https://img.youtube.com/vi/2cLXpOMFGdc/0.jpg" alt="Watch the tutorial" width="240"/>
</a>

## How it works

**4scraper** was developed as an exercise in getting my hands dirty with Go, so there's nothing wild going on behind the scenes. Still, this is how it works, in case anyone's interested.

1. Once a thread URL is provided, the `board` and `threadId` are extracted for later, and an indeterminate [ProgressBar](https://github.com/schollz/progressbar/) is initiated and shown on screen
2. Next, [Colly](https://github.com/gocolly/colly) is used to scrape and retrieve all file URLs and filenames in the thread; files are found by looking for `.filetext > a`
3. All files found are stored in a `[]DownloadableFile slice` for later and the `ProgressBar` is updated to reflect the total files to download
4. We create directories in this structure: `downloads/<board>/<threadId>` to hold all downloaded files
5. For each `DownloadableFile` found we first check if the filename already exists (and append a random number to the filename if it does) then we download it from 4chan

## Found a bug? Have suggestions?

Feel free to use the Issues tab above (or [click here](https://github.com/criticalsession/4scraper/issues)) if you've found bugs, have problems running **4scraper**, have suggestions for improvements or general tips on how I can make the Go code better.

## Things I want to add

- [ ] Optional config file for basic settings (folder organization, types of files to download)
- [ ] Args to skip input with zero feedback to allow for automation

## Like 4scraper?

If you're feeling generous, buy me a beer! - https://www.buymeacoffee.com/criticalsession 🍺❤️