package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://www.xiaoyuzhoufm.com/episode/69392768281939cce65925d3"

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find all script tags (often JSON data is here)
	fmt.Println("=== Script Tags ===")
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		text := s.Text()
		if src != "" {
			fmt.Printf("Script %d: src=%s\n", i, src)
		}
		if len(text) > 0 && len(text) < 500 {
			fmt.Printf("Script %d content: %s\n", i, text[:min(200, len(text))])
		}
	})

	// Look for meta tags
	fmt.Println("\n=== Meta Tags ===")
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		property, _ := s.Attr("property")
		name, _ := s.Attr("name")
		content, _ := s.Attr("content")
		if property != "" || name != "" {
			fmt.Printf("Meta %d: property=%s name=%s content=%s\n", i, property, name, content[:min(100, len(content))])
		}
	})

	// Look for JSON-LD
	fmt.Println("\n=== JSON-LD ===")
	doc.Find("script[type='application/ld+json']").Each(func(i int, s *goquery.Selection) {
		fmt.Printf("JSON-LD %d: %s\n", i, s.Text())
	})

	// Find all audio/source tags
	fmt.Println("\n=== Audio/Source Tags ===")
	doc.Find("audio").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		fmt.Printf("Audio %d: src=%s\n", i, src)
	})
	doc.Find("source").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		fmt.Printf("Source %d: src=%s\n", i, src)
	})

	// Look for data attributes
	fmt.Println("\n=== Elements with data-audio ===")
	doc.Find("[data-audio]").Each(func(i int, s *goquery.Selection) {
		audio, _ := s.Attr("data-audio")
		fmt.Printf("Element %d: data-audio=%s\n", i, audio)
	})

	fmt.Println("\n=== Done ===")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
