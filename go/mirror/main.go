package main

import (
	"github.com/myprivatealaska/bradfield-CSI-prep/go/mirror/crawler"
)

func main() {
	c := crawler.New("https://golang.org")
	c.Crawl()
}
