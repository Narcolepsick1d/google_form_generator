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
func IsProb(prop string) bool {
	regex := `^\d+(,\d+)*$`
	re := regexp.MustCompile(regex)
	return re.Match([]byte(prop))
}
func GetLabel(htmls []string) ([]model.Label, []string) {
	resp := make([]model.Label, 0)
	firstStr := `data-params="%.@.[`
	lastStr := `<div jscontroller`
	htmlsResp := make([]string, 0, 10)
	for _, j := range htmls {
		var entry model.Label
		entries := make([]string, 0, 10)
		firstIndex := strings.Index(j, firstStr)
		lastIndex := strings.Index(j, lastStr)
		matcherStr := j[firstIndex+len(firstStr) : lastIndex]
		//fmt.Println(matcherStr)
		htmlsResp = append(htmlsResp, matcherStr)
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
		for _, e := range entries {
			entry.Entry = e
			entry.Name = n[1]
			resp = append(resp, entry)
		}
	}
	return resp, htmlsResp
}
func GetChoices(htmls []string) [][]string {
	regex := `\[&#34;([^&#]*)&#34;,null`
	re := regexp.MustCompile(regex)

	samplex := make([][]string, 0, 10)
	for _, v := range htmls {
		sample := make([]string, 0, 10)
		if re.Match([]byte(v)) {
			m := re.FindAllStringSubmatch(v, -1)
			for i := 0; i < len(m); i++ {
				if m[i][1] != "" {
					sample = append(sample, m[i][1])
				}
			}
			samplex = append(samplex, sample)
		} else { //для матриц и ебаный строчки * - * - * - * на сколько вы пидораз от 1 до 10
			hash := make(map[string]int)
			regexp1 := `\[&#34;([^&#]*)&#34;\]`
			rep := regexp.MustCompile(regexp1)
			n := rep.FindAllStringSubmatch(v, -1)
			for i := 0; i < len(n); i++ {
				if n[i][1] != "" {
					sample = append(sample, n[i][1])
				}
			}
			dig := strings.Join(sample, "")
			digre := `^\d+`
			dick := regexp.MustCompile(digre)
			if dick.Match([]byte(dig)) {
				samplex = append(samplex, sample)
			} else {
				var res []string
				for _, vi := range sample {
					hash[vi]++
				}
				for _, item := range sample {
					if hash[item] > 1 {
						res = append(res, item)
					}
				}
				//for i, r := range res {
				//	div:=
				//}
				div := res[0]
				count := 0
				for _, r := range res {
					if strings.Contains(r, div) {
						count++
					}
				}
				for d := 0; d < count; d++ {
					samplex = append(samplex, res[0:count])
				}

			}
		}
	}
	return samplex
}
