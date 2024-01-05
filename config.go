package main

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	BoardDir            bool
	ThreadDir           bool
	UseOriginalFilename bool
	Extensions          []string
	ParallelDownload    bool
}

func ReadConfig() Config {
	var c = Config{ // set defaults
		BoardDir:            true,
		ThreadDir:           true,
		UseOriginalFilename: true,
		ParallelDownload:    true,
		Extensions: []string{"jpeg", "jpg", "png", "gif", "webp", "bpm",
			"tiff", "svg", "mp4", "mov", "avi", "webm", "flv"},
	}

	const configFile string = "4scraper.config"

	if FileExists("", configFile) {
		file, err := os.Open(configFile)
		if err != nil {
			return c
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			// skip if too short, starts with # or has no equals sign
			if strings.HasPrefix(line, "#") || len(line) < 5 || !strings.Contains(line, "=") {
				continue
			}

			parts := strings.Split(line, "=")
			key := strings.ToLower(strings.TrimSpace(parts[0]))
			val := strings.Replace(
				strings.ToLower(strings.TrimSpace(parts[1])),
				"\"", "", -1)

			switch key {
			case "paralleldownload":
				c.ParallelDownload = val == "true"
			case "boarddir":
				c.BoardDir = val == "true"
			case "threaddir":
				c.ThreadDir = val == "true"
			case "useoriginalfilename":
				c.UseOriginalFilename = val == "true"
			case "extensions":
				if len(val) == 0 { // empty
					c.Extensions = []string{}
				} else if !strings.Contains(val, ",") { // single extension
					c.Extensions = []string{val}
				} else {
					extensions := strings.Split(val, ",")
					c.Extensions = []string{}
					for _, ext := range extensions {
						if len(strings.TrimSpace(ext)) > 0 {
							c.Extensions = append(c.Extensions, strings.TrimSpace(ext))
						}
					}
				}
			}
		}
	}

	return c
}
