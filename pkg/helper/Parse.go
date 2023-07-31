package helper

import (
	"github.com/PuerkitoBio/goquery"
	"google-gen/internal/model"
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
	return strings.Contains(url, "google.com/forms/") && strings.Contains(url, "viewform") || strings.Contains(url, "https://forms.gle/")
}
func GetLabel(htmls []string) []model.Label {
	resp := make([]model.Label, 0)
	firstStr := `data-params="%.@.[`
	for _, j := range htmls {
		var entry model.Label
		entries := make([]string, 0)
		firstIndex := strings.Index(j, firstStr)
		lastIndex := strings.Index(j, "<div jscontroller")
		matcherStr := j[firstIndex+len(firstStr) : lastIndex]
		regex := `\[(\d+),`
		re := regexp.MustCompile(regex)
		m := re.FindAllStringSubmatch(matcherStr, -1)

		for i := 0; i < len(m); i++ {
			if m[i][1] != "" {
				entries = append(entries, m[i][1])
			}
		}
		regex = `\d+,&#34;(.+?)&#34;,`
		rep := regexp.MustCompile(regex)
		n := rep.FindStringSubmatch(matcherStr)
		log.Println("dsdasd", entries)
		for _, e := range entries {
			log.Print("entry :", e)
			entry.Entry = e
			entry.Name = n[1]
			resp = append(resp, entry)
		}

	}

	return resp
}
