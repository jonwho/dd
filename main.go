package main

import (
	"container/list"
	"flag"
	"fmt"
	// "io/ioutil"
	"log"
	"regexp"

	"github.com/gocolly/colly"
)

const (
	wsjSource = "https://quotes.wsj.com/"
)

var (
	ticker          string
	nLinks          int
	links           *list.List
	visitCount      int
	visited         map[string]bool
	visitWhitelist  = regexp.MustCompile(`(wsj.com|barrons.com|marketwatch.com)`)
	scrapeWhitelist = regexp.MustCompile(`(wsj.com/articles|barrons.com/articles|marketwatch.com/story)`)
)

func init() {
	log.Println("init() starting...")
	flag.StringVar(&ticker, "t", "spy", "Ticker")
	flag.IntVar(&nLinks, "n", 10, "Number of links to visit")
	flag.Parse()

	visited = map[string]bool{}
}

func main() {
	log.Println("Running with ticker", ticker)
	log.Println("Running with link limit", nLinks)
	url := fmt.Sprintf("%s%s", wsjSource, ticker)
	log.Println("First URL", url)

	// new doubly-linked list; use as a queue
	links := list.New()
	links.PushBack(url)

	for link := links.Front(); link != nil && visitCount < nLinks; link = link.Next() {
		visitCount++
		moreLinks := visit(link)
		for _, moreLink := range moreLinks {
			log.Println("Adding link", moreLink)
			links.PushBack(moreLink)
		}
		visited[link.Value.(string)] = true
	}
	return
}

func visit(link *list.Element) []string {
	moreLinks := []string{}
	log.Println("Visiting link", link.Value)
	c := colly.NewCollector()

	c.OnHTML("body", func(e *colly.HTMLElement) {
		log.Println("HTML DOC HERE")
		log.Println(e.Text)
	})

	c.OnHTML("body a", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if !visitWhitelist.MatchString(href) || visited[href] {
			log.Println("WONT VISIT OR SEEN BEFORE. SKIPPING ---", href)
			return
		}

		moreLinks = append(moreLinks, href)
		visited[href] = true
	})

	c.Visit(link.Value.(string))

	return moreLinks
}
