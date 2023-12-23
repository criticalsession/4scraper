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
}

func ReadConfig() Config {
	var config = Config{ // set defaults
		BoardDir:            true,
		ThreadDir:           true,
		UseOriginalFilename: true,
		Extensions: []string{"jpeg", "jpg", "png", "gif", "webp", "bpm",
			"tiff", "svg", "mp4", "mov", "avi", "webm", "flv"},
	}

	const configFile string = "4scraper.config"

	if FileExists("", configFile) {
		file, err := os.Open(configFile)
		if err != nil {
			return config
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			// skip if too short, starts with # or has no equals sign
			if strings.HasPrefix(line, "#") || len(line) < 5 || strings.Index(line, "=") == -1 {
				continue
			}

			parts := strings.Split(line, "=")
			key := strings.ToLower(strings.TrimSpace(parts[0]))
			value := strings.Replace(
				strings.ToLower(strings.TrimSpace(parts[1])),
				"\"", "", -1)

			switch key {
			case "boarddir":
				config.BoardDir = value == "true"
			case "threaddir":
				config.ThreadDir = value == "true"
			case "useoriginalfilename":
				config.UseOriginalFilename = value == "true"
			case "extensions":
				if len(value) == 0 { // empty
					config.Extensions = []string{}
				} else if strings.Index(value, ",") == -1 { // single extension
					config.Extensions = []string{value}
				} else {
					extensions := strings.Split(value, ",")
					config.Extensions = []string{}
					for _, ext := range extensions {
						if len(strings.TrimSpace(ext)) > 0 {
							config.Extensions = append(config.Extensions, strings.TrimSpace(ext))
						}
					}
				}
			}
		}
	}

	return config
}
