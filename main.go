package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Target struct {
	Endpoint string
	Element  string
	Keyword  string
}

type Result struct {
	Endpoint string
	Content  string
	Error    error
}

func scrape(target Target) (content string, err error) {
	resp, err := http.Get(target.Endpoint)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	var episodes []string
	doc.Find(target.Element).Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.Contains(strings.ToLower(text), target.Keyword) {
			link, exists := s.Find("a").First().Attr("href")
			if !exists {
				episodes = append(episodes, fmt.Sprintf("Error: Link not found for episode: %s", text))
				return
			}
			episodes = append(episodes, fmt.Sprintf("%s - %s", text, link))
		}
	})

	jsonContent, err := json.Marshal(episodes)
	if err != nil {
		return "", err
	}
	return string(jsonContent), nil
}

func main() {
	targets := []Target{
		{"https://kubernetespodcast.com/", "div.episode h3", "istio"},
		{"https://www.heroku.com/coderish", "div.episode h3", "engineering"},
	}

	var wg sync.WaitGroup
	resChan := make(chan Result, len(targets))

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			for _, target := range targets {
				wg.Add(1)
				go func(t Target) {
					defer wg.Done()
					content, err := scrape(t)
					resChan <- Result{Endpoint: t.Endpoint, Content: content, Error: err}
				}(target)
			}
		}
	}()

	go func() {
		for result := range resChan {
			if result.Error != nil {
				fmt.Printf("Error scraping %s: %v\n", result.Endpoint, result.Error)
			} else {
				fmt.Printf("Content from %s: %s\n", result.Endpoint, result.Content)
			}
		}
	}()

	wg.Wait()
	close(resChan)
}
