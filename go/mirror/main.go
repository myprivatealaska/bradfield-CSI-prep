package main

import (
	"sync"

	"github.com/myprivatealaska/bradfield-CSI-prep/go/mirror/crawler"
)

func main() {
	wg := sync.WaitGroup{}
	c := crawler.New("https://golang.org")
	wg.Add(1)
	go c.Crawl(c.EntryUrl, 1, &wg)
	wg.Wait()
}
