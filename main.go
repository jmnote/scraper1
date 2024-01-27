package main

import (
	"fmt"
	"io"
	"net/http"
)

func scrape(endpoint string) (content string, err error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
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
