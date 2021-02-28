package crawler

import (
	io "io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"mvdan.cc/xurls/v2"
)

type crawler struct {
	EntryUrl       string
	MaxDepth       int
	CurDepth       int
	Visited        map[string]bool
	NeedProcessing chan string
	//Errors []error
}

func New(entryUrl string) crawler {
	c := crawler{
		EntryUrl: entryUrl,
		MaxDepth: 2,
		Visited:  map[string]bool{},
		//NeedProcessing: make(chan string, 4),
	}
	return c
}

func (c crawler) Crawl() {
	// Crawl until url is nil and all processPage runs have finished
	//for ; ; {
	//	url := <- c.NeedProcessing
	//	if !c.Visited[url] {
	//		c.Visited[url] = true
	//		go c.processPage(url)
	//	}
	//}
	c.processPage(c.EntryUrl)
}

func (c crawler) processPage(url string) {
	c.CurDepth++
	// outside: check if hasn't been visited
	client := http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		log.Println(errors.Wrap(err, "processPage - can't create http client"))
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(errors.Wrap(err, "processPage - can't read response"))
	}
	strBody := string(body)
	rxStrict := xurls.Strict()
	urls := rxStrict.FindAllString(strBody, 15)
	for u := range urls {
		println(u)
	}

	//if c.CurDepth < c.MaxDepth {
	//	c.NeedProcessing <- "next_url"
	//}
}
