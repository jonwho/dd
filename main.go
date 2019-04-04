package main

import (
	"container/list"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"sync"
	"sync/atomic"

	"github.com/gocolly/colly"
	"github.com/google/uuid"
)

const (
	wsjSource        = "https://quotes.wsj.com/"
	wsjAAPLSource    = "https://quotes.wsj.com/aapl"
	wsjMSFTSource    = "https://quotes.wsj.com/msft"
	wsjTSLASource    = "https://quotes.wsj.com/tsla"
	wsjFBSource      = "https://quotes.wsj.com/fb"
	writePermissions = 0644
)

var (
	ticker          string
	nLinks          uint64
	visitCount      uint64
	visited         sync.Map
	visitWhitelist  = regexp.MustCompile(`(wsj.com|barrons.com|marketwatch.com)`)
	scrapeWhitelist = regexp.MustCompile(`(wsj.com\/articles|barrons.com\/articles|marketwatch.com\/story)`)
	wg              sync.WaitGroup
)

func init() {
	log.Println("init() starting...")
	flag.StringVar(&ticker, "t", "aapl", "Ticker")
	flag.Uint64Var(&nLinks, "n", 10, "Number of links to visit")
	flag.Parse()
}

func main() {
	wg.Add(4)

	go crawl(wsjAAPLSource)
	go crawl(wsjMSFTSource)
	go crawl(wsjTSLASource)
	go crawl(wsjFBSource)

	wg.Wait()
}

func crawl(url string) {
	defer wg.Done()

	links := list.New()
	links.PushBack(url)

	for link := links.Front(); link != nil && atomic.LoadUint64(&visitCount) < nLinks; link = link.Next() {
		moreLinks := visit(link)
		for _, moreLink := range moreLinks {
			log.Println("Adding link", moreLink)
			links.PushBack(moreLink)
		}
		visited.Store(link.Value.(string), true)
	}
}

func visit(link *list.Element) []string {
	moreLinks := []string{}
	log.Println("Visiting link", link.Value)
	c := colly.NewCollector(
		colly.URLFilters(
			visitWhitelist,
			scrapeWhitelist,
		),
	)

	c.OnHTML("body", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()
		uri := fmt.Sprintf("tmp/html_%s", uuid.New().String())
		if scrapeWhitelist.MatchString(url) {
			atomic.AddUint64(&visitCount, 1)
			content := append([]byte(url+"\n\n\n"), []byte(e.Text)...)
			log.Println("Writing to file... ", uri)
			err := ioutil.WriteFile(uri, content, writePermissions)
			if err != nil {
				log.Fatalln(err)
			}
		}
	})

	c.OnHTML("body a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if _, ok := visited.Load(href); ok || !visitWhitelist.MatchString(href) {
			visited.Store(href, true)
			log.Println("WONT VISIT OR SEEN BEFORE. SKIPPING ---", href)
			return
		}

		moreLinks = append(moreLinks, href)
		visited.Store(href, true)
	})

	c.Visit(link.Value.(string))

	return moreLinks
}
