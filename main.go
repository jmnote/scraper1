package main

import (
	"fmt"
	"io"
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

	// HTML 문서 읽기
	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// goquery 문서 생성
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
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
				err = fmt.Errorf("link not found for episode: %s", text)
				return
			}
			episodes = append(episodes, fmt.Sprintf("%s - %s", text, link))
		}
	})

	if err != nil {
		return "", err
	}

	return fmt.Sprint(episodes), nil
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
