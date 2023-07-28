package helper

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

func ExampleScrape(url string) []string {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	result := make([]string, 0)
	doc.Find("div.Qr7Oae").Each(func(i int, s *goquery.Selection) {
		s.Find("div").Html()
		html, err := s.Html()
		if err != nil {
			return
		}
		result = append(result, html)
		//fmt.Printf("Review %d: %s\n", i, title)
	})
	return result
}
