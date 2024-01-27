package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func scrape(endpoint string) (content string, err error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// goquery 문서 생성
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// 원하는 요소 선택 및 처리
	var episodes []string
	doc.Find("div.episode h3").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.Contains(strings.ToLower(text), "istio") {
			link, exists := s.Find("a").First().Attr("href")
			if !exists {
				episodes = append(episodes, fmt.Sprintf("Error: Link not found for episode: %s", text))
				return
			}
			episodes = append(episodes, fmt.Sprintf("%s - %s", text, link))
		}
	})

	// JSON 형식으로 변환
	jsonContent, err := json.Marshal(episodes)
	if err != nil {
		return "", err
	}
	return string(jsonContent), nil
}

func main() {
	endpoints := []string{"https://kubernetespodcast.com/"}
	for _, endpoint := range endpoints {
		content, err := scrape(endpoint)
		if err != nil {
			fmt.Println("Error scraping:", err)
			continue
		}
		fmt.Println(content)
	}
}

// https://kubernetespodcast.com/ istio
// https://heroku.com/ engineering
