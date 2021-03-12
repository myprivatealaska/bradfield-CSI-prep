package crawler

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

// crawler struct is private. This is because I don't want the consumer of the package to be able to
// create instances of crawler directly because they might create an empty instance and it will behave
// unexpectedly. I am exposing the constructor - New method - instead
type crawler struct {
	EntryUrl    string
	MaxDepth    int
	VisitedLock sync.Mutex
	Visited     map[string]bool
	Tokens      chan struct{}
}

// This is a constructor. In go there are no classes and constructors per se, but this is how
// you mimic it
func New(entryUrl string) crawler {
	c := crawler{
		EntryUrl: entryUrl,
		MaxDepth: 3,
		// I was getting deadlocks when trying to use sync.Map - Need to figure out why
		VisitedLock: sync.Mutex{},
		Visited:     make(map[string]bool),
		// tokens is a counting semaphore used to
		// enforce a limit of 20 concurrent requests.
		Tokens: make(chan struct{}, 20),
	}
	return c
}

// Crawl will be called recursively until c.MaxDepth is reached and all processUrl runs have finished
func (c crawler) Crawl(url string, curDepth int, wg *sync.WaitGroup) {
	defer wg.Done()

	c.Tokens <- struct{}{} // acquire a token
	foundUrls, err := c.processUrl(url)
	<-c.Tokens //release token

	if err != nil {
		log.Println(fmt.Sprintf("failed to process %v: %e", url, err))
	}

	if curDepth >= c.MaxDepth {
		return
	}

	// Let's put visited urls in a map to avoid visiting them again in the future
	for _, link := range foundUrls {
		c.VisitedLock.Lock()
		if c.Visited[link] {
			c.VisitedLock.Unlock()
			continue
		} else {
			c.Visited[link] = true
			c.VisitedLock.Unlock()
			wg.Add(1)
			// queue for processing
			log.Println(fmt.Sprintf("queue for processing: %v", link))
			go c.Crawl(link, curDepth+1, wg)
		}
	}
}

// processUrl fetches the page by url, processes it, and stores a local copy of the page
func (c crawler) processUrl(url string) ([]string, error) {
	client := http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		log.Println(errors.Wrap(err, "processUrl - can't create http client"))
	}

	defer resp.Body.Close()
	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		log.Println(fmt.Sprintf("not html, skipping %v. %v", url, resp.Header.Get("Content-Type")))
		return []string{}, nil
	}
	foundUrls, pageBytes, err := processPage(url, resp.Body)
	if err != nil {
		log.Println(fmt.Sprintf("can't processUrl %v. error: %e", resp.Request.URL, err))
		return []string{}, err
	}
	if err := storePage(resp.Request.URL, pageBytes); err != nil {
		log.Println(fmt.Sprintf("can't store the page %v. error: %e", resp.Request.URL, err))
		return foundUrls, err
	}

	return foundUrls, nil
}

// processPage parses urls we care about from the page, overwrites them for storing locally,
// and returns an array of unique URLs and a bytes.Buffer ready to be stored locally
func processPage(rawUrl string, body io.Reader) (parsedUrls []string, updatedPage io.Reader, err error) {
	// Let's parse our input Url as a url to use the functionality of the url package
	u, err := url.Parse(rawUrl)
	if err != nil {
		err = errors.Wrap(err, "parsePage - can't parse rawURL")
	}
	// rootNode is now a pointer to the root node of the HTML document. It's essentially the head
	// of the Linked List
	rootNode, err := html.Parse(body)
	if err != nil {
		err = errors.Wrap(err, "parsePage")
	}
	nodes := linkNodes(rootNode)
	parsedUrlsMap := linkURLs(nodes, u) // Extract unique URL from the page
	for k, _ := range parsedUrlsMap {
		parsedUrls = append(parsedUrls, k)
	}
	rewriteLocalLinks(nodes, u)
	b := bytes.Buffer{}
	if errRend := html.Render(&b, rootNode); err != nil {
		err = errors.Wrap(errRend, "parsePage - can't build html string back")
		return
	}
	updatedPage = &b
	return
}

// storePage saves pageBody to the local fs. It creates all dirs as necessary.
func storePage(pageUrl *url.URL, pageBody io.Reader) error {
	fileName := filepath.Join(pageUrl.Host, pageUrl.Path)
	if filepath.Ext(pageUrl.Path) == "" {
		fileName = filepath.Join(pageUrl.Host, pageUrl.Path, "index.html")
	}
	err := os.MkdirAll(filepath.Dir(fileName), 0777)
	if err != nil {
		return errors.Wrap(err, "storePage - failed to create dir(s)")
	}
	file, err := os.Create(fileName)
	if err != nil {
		return errors.Wrap(err, "storePage - failed to create file")
	}
	if _, err = io.Copy(file, pageBody); err != nil {
		return errors.Wrap(err, "storePage - failed to copy pageBody to file")
	}
	// Check for delayed write errors, as mentioned at the end of section 5.8.
	err = file.Close()
	if err != nil {
		return errors.Wrap(err, "storePage - delayed write errors")
	}
	return nil
}

// forEachNode calls the functions pre(x) and post(x) for each node
// x in the tree rooted at n. Both functions are optional.
// pre is called before the children are visited (preorder) and
// post is called after (postorder). (5.5 Function Values)
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

// linkNodes parses out all nodes that are links given the root node.
// we will re-use the result of it in both linkURLs and rewriteLocalLinks
// that's why it's a separate function
func linkNodes(n *html.Node) []*html.Node {
	var links []*html.Node
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			links = append(links, n)
		}
	}
	forEachNode(n, visitNode, nil)
	return links
}

// linkURLs returns a map of URLs filtered so that it's either links with the same domain
// as c.EntryUrl or relative URLs
func linkURLs(linkNodes []*html.Node, base *url.URL) (parsedUrls map[string]bool) {
	parsedUrls = map[string]bool{}
	for _, n := range linkNodes {
		for _, a := range n.Attr {
			if a.Key == "href" {
				// no need to use regex, Parse is smart!
				link, err := base.Parse(a.Val)
				// ignore bad and non-local URLs
				if err != nil || link.Host != base.Host {
					continue
				}
				urlString := link.String()
				if !parsedUrls[urlString] {
					parsedUrls[urlString] = true
				}
			}
		}
	}
	return
}

// rewriteLocalLinks rewrites local links to be relative and links without
// extensions to point to index.html. Since linkNodes is a slice of pointers,
// we're essentially modifying what these pointers point to.
// after rewriteLocalLinks does its job, we can serialize the rootNode back to a HTML string
func rewriteLocalLinks(linkNodes []*html.Node, base *url.URL) {
	for _, n := range linkNodes {
		for i, a := range n.Attr {
			if a.Key == "href" {
				link, err := base.Parse(a.Val)
				if err != nil || link.Host != base.Host {
					continue // ignore bad and non-local URLs
				}
				// Make these links relative
				link.Scheme = ""
				link.Host = ""
				link.User = nil
				a.Val = link.String()
				// Change the href attr of the link
				n.Attr[i] = a
			}
		}
	}
}
