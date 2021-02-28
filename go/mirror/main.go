package main

import (
	"crawler"
)

func main() {
	c := crawler.New("https://golang.org")
	c.Crawl()
}
