package search

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type Page struct {
	Page    int
	Threads []Thread
}

type Thread struct {
	No     int
	Sticky int
	Com    string
	Sub    string
}

func (t *Thread) getCombinedQueryText() string {
	return strings.ToLower(t.Com + " " + t.Sub)
}

func FindInBoard(board string, keywords []string) ([]Thread, error) {
	threads := []Thread{}

	pages, err := parseBoard(board)
	if err != nil {
		return nil, err
	}

	for _, p := range pages {
		for _, t := range p.Threads {
			if t.Sticky == 1 || t.getCombinedQueryText() == "" {
				continue
			}

			match := true
			for _, k := range keywords {
				if !strings.Contains(t.getCombinedQueryText(), strings.ToLower(k)) {
					match = false
					break
				}
			}

			if match {
				threads = append(threads, t)
			}
		}
	}

	return threads, nil
}

func parseBoard(board string) ([]Page, error) {
	resp, err := http.Get("https://a.4cdn.org/" + board + "/catalog.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result []Page
	err = json.Unmarshal(body, &result)

	return result, err
}
