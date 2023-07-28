package helper

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"regexp"
	"strings"
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
func IsGoogleFormsLink(url string) bool {
	// Проверяем, содержит ли ссылка подстроку "google.com/forms/"
	// и "viewform"
	return strings.Contains(url, "google.com/forms/") && strings.Contains(url, "viewform")
}
func GetEntry(htmls []string) []string {
	resp := make([]string, 0)
	firstStr := `data-params="%.@.`
	for _, j := range htmls {
		firstIndex := strings.Index(j, firstStr)
		lastIndex := strings.Index(j, "<div jscontroller")
		matcherStr := j[firstIndex+len(firstStr) : lastIndex]
		regex := `\[\[(\d+),`
		re := regexp.MustCompile(regex)
		m := re.FindStringSubmatch(matcherStr)
		if len(matcherStr) > 1 {
			resp = append(resp, m[1])
		}
	}
	return resp
}
