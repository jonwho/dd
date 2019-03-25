package dd

import (
	"log"
	"strings"

	"github.com/gocolly/colly"
)

const (
	wsjSource = "https://quotes.wsj.com/"
)

// WSJ get news for ticker from Wall Street Journal
func WSJ(ticker string) ([]string, error) {
	c := colly.NewCollector()
	var content []string
	var err error

	// Register callback for every matching selector
	c.OnHTML("span.headline > a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		content = append(content, marketWatchArticle(link))
	})

	c.Visit(wsjSource + ticker)
	log.Println(content[0])
	return content, err
}

// Return article text from MarketWatch
func marketWatchArticle(url string) string {
	sb := strings.Builder{}
	c := colly.NewCollector()

	c.OnHTML("article#article p", func(e *colly.HTMLElement) {
		sb.WriteString(e.Text)
	})

	c.Visit(url)

	return sb.String()
}

// @todo: missing domains to pull articles from
// - wsj
// - barrons
